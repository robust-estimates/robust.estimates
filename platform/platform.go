package platform

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	core "robust/platform/core"
)

type Claims = core.Claims

type Error = core.Error

type Command = core.Command

var E = core.E

var Invalid = core.Invalid

var Conflict = core.Conflict

var NotFound = core.NotFound

var Forbidden = core.Forbidden

var ID = core.ID

var Now = core.Now

var WithIdentity = core.WithIdentity

var Identity = core.Identity

var Require = core.Require

var Hash = core.Hash

var GetCommand = core.GetCommand

var CheckVersion = core.CheckVersion

func Env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}

	return fallback
}

func Secret() string {
	s := os.Getenv("JWT_SECRET")
	if len(s) < 32 {
		panic("JWT_SECRET must contain at least 32 bytes")
	}

	return s
}

type tokenClaims struct {
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func Token(secret, user, name, role string) (string, error) {
	c := tokenClaims{Role: role, Name: name, RegisteredClaims: jwt.RegisteredClaims{Subject: user, Issuer: "robust", Audience: jwt.ClaimStrings{"credit"}, IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secret))
}

func Parse(secret, raw string) (Claims, error) {
	var c tokenClaims
	_, e := jwt.ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) { return []byte(secret), nil }, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("robust"), jwt.WithAudience("credit"), jwt.WithExpirationRequired())
	if e != nil || c.Subject == "" || (c.Role != "CLIENT" && c.Role != "EMPLOYEE" && c.Role != "SYSTEM") {
		return Claims{}, E("UNAUTHORIZED", "Требуется вход", 401)
	}

	return Claims{Subject: c.Subject, Name: c.Name, Role: c.Role}, nil
}

func Incoming(ctx context.Context, k string) string {
	m, _ := metadata.FromIncomingContext(ctx)
	v := m.Get(k)
	if len(v) > 0 {
		return v[0]
	}

	return ""
}

func RPCError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	var e *Error
	if errors.As(err, &e) {
		c := codes.Internal
		switch e.Status {
		case 400:
			c = codes.InvalidArgument
		case 401:
			c = codes.Unauthenticated
		case 403:
			c = codes.PermissionDenied
		case 404:
			c = codes.NotFound
		case 409:
			c = codes.Aborted
		case 413:
			c = codes.ResourceExhausted
		case 415, 422:
			c = codes.FailedPrecondition
		}

		b, _ := json.Marshal(e)

		return status.Error(c, string(b))
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, "Превышено время ожидания")
	}

	slog.Error("operation failed", "error", err)
	return status.Error(codes.Internal, "Внутренняя ошибка сервиса")
}

func HTTPError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	s, ok := status.FromError(err)
	if !ok {
		return &Error{Code: "INTERNAL", Message: "Внутренняя ошибка", Status: 500}
	}

	var mapped Error
	if json.Unmarshal([]byte(s.Message()), &mapped) == nil && mapped.Code != "" {
		switch s.Code() {
		case codes.InvalidArgument:
			mapped.Status = 400
		case codes.Unauthenticated:
			mapped.Status = 401
		case codes.PermissionDenied:
			mapped.Status = 403
		case codes.NotFound:
			mapped.Status = 404
		case codes.Aborted:
			mapped.Status = 409
		case codes.ResourceExhausted:
			mapped.Status = 413
		case codes.FailedPrecondition:
			mapped.Status = 422
			if mapped.Code == "UNSUPPORTED_MEDIA_TYPE" {
				mapped.Status = 415
			}
		default:
			mapped.Status = 500
		}

		return &mapped
	}

	h := 500
	if s.Code() == codes.Unavailable {
		h = 503
	}

	if s.Code() == codes.DeadlineExceeded {
		h = 504
	}

	return &Error{Code: "UPSTREAM_ERROR", Message: s.Message(), Status: h}
}

func WriteError(w http.ResponseWriter, err error, trace string) {
	e := HTTPError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(map[string]any{"code": e.Code, "message": e.Message, "traceId": trace})
}

func Dial(address string) *grpc.ClientConn {
	c, e := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if e != nil {
		panic(e)
	}

	return c
}

func Outgoing(ctx context.Context) context.Context {
	m, _ := metadata.FromIncomingContext(ctx)

	return metadata.NewOutgoingContext(ctx, m.Copy())
}

func SystemContext(secret string) (context.Context, context.CancelFunc) {
	t, e := Token(secret, "worker", "worker", "SYSTEM")
	if e != nil {
		panic(e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+t)), cancel
}

func authenticate(ctx context.Context, secret string) (context.Context, error) {
	raw := strings.TrimPrefix(Incoming(ctx, "authorization"), "Bearer ")
	c, e := Parse(secret, raw)
	if e != nil {
		return ctx, RPCError(e)
	}

	return core.WithHeaders(WithIdentity(ctx, c), Incoming(ctx, "idempotency-key"), Incoming(ctx, "if-match")), nil
}

func UnaryAuth(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
		if strings.HasPrefix(info.FullMethod, "/grpc.health.") {
			return h(ctx, req)
		}

		ctx, e := authenticate(ctx, secret)
		if e != nil {
			return nil, e
		}

		v, e := h(ctx, req)

		return v, RPCError(e)
	}
}

type authStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authStream) Context() context.Context { return s.ctx }

func StreamAuth(secret string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		ctx, e := authenticate(stream.Context(), secret)
		if e != nil {
			return e
		}

		return RPCError(h(srv, &authStream{stream, ctx}))
	}
}

func NewGRPC(secret string) *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryAuth(secret)), grpc.StreamInterceptor(StreamAuth(secret)), grpc.MaxRecvMsgSize(1024*1024), grpc.MaxSendMsgSize(1024*1024))
	h := health.NewServer()
	h.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, h)
	return s
}

func Serve(s *grpc.Server, address string, cleanup func()) {
	l, e := net.Listen("tcp", address)
	if e != nil {
		panic(e)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		done := make(chan struct{})
		go func() { s.GracefulStop(); close(done) }()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			s.Stop()
		}

		cleanup()
	}()
	slog.Info("grpc listening", "address", address)
	if e = s.Serve(l); e != nil {
		panic(e)
	}
}

func OpenDB(ctx context.Context, dsn string) *pgxpool.Pool {
	p, e := pgxpool.New(ctx, dsn)
	if e != nil {
		panic(e)
	}

	if e = p.Ping(ctx); e != nil {
		panic(e)
	}

	return p
}

func ToProto(v any, p proto.Message) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}

	return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(b, p)
}

func ProtoJSON(p proto.Message) ([]byte, error) {
	b, e := protojson.MarshalOptions{EmitDefaultValues: true}.Marshal(p)
	if e != nil {
		return nil, e
	}

	var v map[string]any
	if e = json.Unmarshal(b, &v); e != nil {
		return nil, e
	}

	normalize(v, p.ProtoReflect().Descriptor())
	return json.Marshal(v)
}

func normalize(v map[string]any, d protoreflect.MessageDescriptor) {
	fields := d.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		x, ok := v[f.JSONName()]
		if !ok {
			continue
		}

		conv := func(value any) any {
			if f.Kind() == protoreflect.Int64Kind || f.Kind() == protoreflect.Uint64Kind {
				if z, ok := value.(string); ok {
					return json.Number(z)
				}
			}

			if f.Kind() == protoreflect.MessageKind {
				if m, ok := value.(map[string]any); ok {
					normalize(m, f.Message())
				}
			}

			return value
		}
		if f.IsMap() {
			if m, ok := x.(map[string]any); ok {
				kind := f.MapValue().Kind()
				if kind == protoreflect.Int64Kind || kind == protoreflect.Uint64Kind {
					for k, y := range m {
						if z, ok := y.(string); ok {
							m[k] = json.Number(z)
						}
					}
				}
			}

			continue
		}

		if f.IsList() {
			if xs, ok := x.([]any); ok {
				for n, y := range xs {
					xs[n] = conv(y)
				}
			}
		} else {
			v[f.JSONName()] = conv(x)
		}
	}
}

func FromProto(p proto.Message, v any) error {
	b, e := ProtoJSON(p)
	if e != nil {
		return e
	}

	return json.Unmarshal(b, v)
}

func PB[T proto.Message](v any, p T) T {
	if e := ToProto(v, p); e != nil {
		panic(e)
	}

	return p
}

func Healthcheck(address string) bool {
	if len(os.Args) < 2 || os.Args[1] != "healthcheck" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := Dial("127.0.0.1:" + strings.TrimPrefix(address, ":"))
	defer c.Close()
	v, e := healthpb.NewHealthClient(c).Check(ctx, &healthpb.HealthCheckRequest{})
	if e != nil || v.Status != healthpb.HealthCheckResponse_SERVING {
		os.Exit(1)
	}

	return true
}
