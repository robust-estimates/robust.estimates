package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	pb "robust/contracts/credit/v1"
	"robust/gateway/internal/usecase"
	"robust/platform"
)

type route func(http.ResponseWriter, *http.Request) error

type Server struct {
	Auth            *usecase.Auth
	Companies       pb.CompanyServiceClient
	Applications    pb.ApplicationServiceClient
	DB              *pgxpool.Pool
	CompanyConn     *grpc.ClientConn
	ApplicationConn *grpc.ClientConn
	Secret          string
	Origin          string
	OpenAPI         string
	Swagger         string
	mu              sync.Mutex
	attempts        map[string][]time.Time
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.register(mux, "POST /api/v1/auth/login", s.login, false)
	s.register(mux, "POST /api/v1/auth/register", s.signup, false)
	s.register(mux, "GET /api/v1/me", s.me, true)
	s.register(mux, "POST /api/v1/companies", s.createCompany, true)
	s.register(mux, "GET /api/v1/companies", s.listCompanies, true)
	s.register(mux, "GET /api/v1/companies/{companyId}", s.getCompany, true)
	s.register(mux, "PATCH /api/v1/companies/{companyId}", s.updateCompany, true)
	s.register(mux, "POST /api/v1/companies/{companyId}/verification", s.verifyCompany, true)
	s.register(mux, "GET /api/v1/products", s.products, true)
	s.register(mux, "GET /api/v1/products/{productId}", s.product, true)
	s.register(mux, "POST /api/v1/recommendations", s.recommend, true)
	s.register(mux, "POST /api/v1/calculations", s.calculate, true)
	s.register(mux, "POST /api/v1/applications", s.createApplication, true)
	s.register(mux, "GET /api/v1/applications", s.listApplications, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}", s.getApplication, true)
	s.register(mux, "PATCH /api/v1/applications/{applicationId}", s.patchApplication, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/validation", s.validation, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/submit", s.submit, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/cancel", s.cancel, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}/decision", s.decision, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/accept", s.accept, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/decline", s.decline, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/issuance", s.issue, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}/history", s.history, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}/information-requests", s.information, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}/documents", s.documents, true)
	s.register(mux, "POST /api/v1/applications/{applicationId}/documents", s.upload, true)
	s.register(mux, "DELETE /api/v1/applications/{applicationId}/documents/{documentId}", s.deleteDocument, true)
	s.register(mux, "GET /api/v1/applications/{applicationId}/documents/{documentId}/download", s.download, true)
	s.register(mux, "GET /api/v1/bank/applications", s.bankApplications, true)
	s.register(mux, "GET /api/v1/bank/applications/{applicationId}/scoring", s.scoring, true)
	s.register(mux, "GET /api/v1/bank/applications/{applicationId}", s.overview, true)
	s.register(mux, "POST /api/v1/bank/applications/{applicationId}/decision", s.decide, true)
	s.register(mux, "GET /api/v1/feed", s.feed, true)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { write(w, 200, map[string]string{"status": "ok"}) })
	s.register(mux, "GET /readyz", s.ready, false)
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, s.OpenAPI)
	})
	mux.HandleFunc("GET /docs/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, s.Swagger) })
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		trace := platform.ID()
		w.Header().Set("X-Request-ID", trace)
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Add("Vary", "Origin")
			if origin == s.Origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, If-Match")
				w.Header().Set("Access-Control-Expose-Headers", "ETag, X-Request-ID, Content-Disposition")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		start := time.Now()
		defer func() {
			if x := recover(); x != nil {
				slog.Error("http panic", "traceId", trace, "error", x)
				platform.WriteError(w, platform.E("INTERNAL", "Внутренняя ошибка", 500), trace)
			}

			slog.Info("http request", "method", r.Method, "route", r.Pattern, "duration", time.Since(start), "traceId", trace)
		}()
		mux.ServeHTTP(w, r)
	})
}

func (s *Server) register(m *http.ServeMux, pattern string, h route, auth bool) {
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		trace := w.Header().Get("X-Request-ID")
		if auth {
			raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			c, e := platform.Parse(s.Secret, raw)
			if e != nil {
				platform.WriteError(w, e, trace)
				return
			}

			if c.Role == "SYSTEM" {
				platform.WriteError(w, platform.Forbidden(), trace)
				return
			}

			ctx := platform.WithIdentity(r.Context(), c)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+raw, "idempotency-key", r.Header.Get("Idempotency-Key"), "if-match", r.Header.Get("If-Match"), "x-request-id", trace))
			r = r.WithContext(ctx)
			if r.Method == "POST" || r.Method == "PATCH" || r.Method == "DELETE" {
				k := r.Header.Get("Idempotency-Key")
				if len(k) < 8 || len(k) > 128 {
					platform.WriteError(w, platform.E("INVALID_KEY", "Нужен Idempotency-Key длиной 8–128", 400), trace)
					return
				}
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		if e := h(w, r.WithContext(ctx)); e != nil {
			platform.WriteError(w, e, trace)
		}
	})
}

func read(r *http.Request, p proto.Message) error {
	b, e := io.ReadAll(io.LimitReader(r.Body, 1024*1024+1))
	if e != nil {
		return platform.E("INVALID_JSON", "Не удалось прочитать JSON", 400)
	}

	if len(b) > 1024*1024 {
		return platform.E("PAYLOAD_TOO_LARGE", "JSON больше 1 МиБ", 413)
	}

	if e = protojson.Unmarshal(b, p); e != nil {
		return platform.E("INVALID_JSON", "Некорректный JSON или неизвестное поле: "+e.Error(), 400)
	}

	return nil
}

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func reply(w http.ResponseWriter, code int, p proto.Message) error {
	b, e := platform.ProtoJSON(p)
	if e != nil {
		return e
	}

	var v map[string]any
	if e = json.Unmarshal(b, &v); e != nil {
		return e
	}

	redact(v)
	if a, ok := p.(*pb.Application); ok {
		w.Header().Set("ETag", fmt.Sprintf("\"%d\"", a.Version))
	}

	if _, ok := p.(*pb.Company); ok {
		v = flatten(v)
	}

	write(w, code, v)
	return nil
}

func flatten(v map[string]any) map[string]any {
	if d, ok := v["data"].(map[string]any); ok {
		for k, x := range d {
			v[k] = x
		}

		delete(v, "data")
	}

	delete(v, "ownerId")
	return v
}

func array(w http.ResponseWriter, p proto.Message) error {
	b, e := platform.ProtoJSON(p)
	if e != nil {
		return e
	}

	var v map[string]any
	if e = json.Unmarshal(b, &v); e != nil {
		return e
	}

	items, ok := v["items"].([]any)
	if !ok {
		items = []any{}
	}

	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			delete(m, "ownerId")
			if _, ok := m["data"]; ok {
				flatten(m)
			}
		}
	}

	write(w, 200, items)
	return nil
}

func id(r *http.Request) *pb.ID { return &pb.ID{Id: r.PathValue("applicationId")} }

func (s *Server) allowAuth(r *http.Request, bucket string, limit int) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	key := bucket + ":" + ip
	s.mu.Lock()
	if s.attempts == nil {
		s.attempts = map[string][]time.Time{}
	}

	now := time.Now()
	if len(s.attempts) > 10000 {
		for k, v := range s.attempts {
			if len(v) == 0 || now.Sub(v[len(v)-1]) > time.Minute {
				delete(s.attempts, k)
			}
		}
	}

	recent := []time.Time{}
	for _, t := range s.attempts[key] {
		if now.Sub(t) < time.Minute {
			recent = append(recent, t)
		}
	}

	if len(recent) >= limit {
		s.mu.Unlock()
		return platform.E("RATE_LIMITED", "Слишком много попыток; повторите через минуту", 429)
	}

	s.attempts[key] = append(recent, now)
	s.mu.Unlock()
	return nil
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) error {
	if e := s.allowAuth(r, "login", 10); e != nil {
		return e
	}

	var in struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if e := readAuthJSON(w, r, &in); e != nil {
		return e
	}

	if len(in.Password) > 72 || in.Password == "" || len(in.Login) > 320 {
		return platform.E("INVALID_CREDENTIALS", "Неверный логин или пароль", 401)
	}

	v, e := s.Auth.Login(r.Context(), in.Login, in.Password)
	if e != nil {
		return e
	}

	write(w, 200, v)
	return nil
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) error {
	c := platform.Identity(r.Context())
	ids := []string{}
	if c.Role == "CLIENT" {
		v, e := s.Companies.ListCompanies(r.Context(), &pb.Empty{})
		if e != nil {
			return e
		}

		for _, company := range v.Items {
			ids = append(ids, company.Id)
		}
	}

	write(w, 200, map[string]any{"id": c.Subject, "displayName": c.Name, "role": c.Role, "companyIds": ids})
	return nil
}

func (s *Server) createCompany(w http.ResponseWriter, r *http.Request) error {
	in := &pb.CompanyInput{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Companies.CreateCompany(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 201, v)
}

func (s *Server) listCompanies(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Companies.ListCompanies(r.Context(), &pb.Empty{})
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) getCompany(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Companies.GetCompany(r.Context(), &pb.ID{Id: r.PathValue("companyId")})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) updateCompany(w http.ResponseWriter, r *http.Request) error {
	in := &pb.CompanyInput{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Companies.UpdateCompany(r.Context(), &pb.CompanyUpdate{Id: r.PathValue("companyId"), Data: in})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) verifyCompany(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Companies.VerifyCompany(r.Context(), &pb.ID{Id: r.PathValue("companyId")})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) products(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Companies.ListProducts(r.Context(), &pb.Empty{})
	if e != nil {
		return e
	}

	filtered := []*pb.Product{}
	q := strings.ToLower(r.URL.Query().Get("search"))
	tag := r.URL.Query().Get("tag")
	for _, p := range v.Items {
		if q != "" && !strings.Contains(strings.ToLower(p.Name+" "+p.Description), q) {
			continue
		}

		found := tag == ""
		for _, t := range p.Tags {
			if t == tag {
				found = true
			}
		}

		if found {
			filtered = append(filtered, p)
		}
	}

	v.Items = filtered

	return array(w, v)
}

func (s *Server) product(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Companies.GetProduct(r.Context(), &pb.ID{Id: r.PathValue("productId")})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) recommend(w http.ResponseWriter, r *http.Request) error {
	in := &pb.RecommendationInput{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Companies.Recommend(r.Context(), in)
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) calculate(w http.ResponseWriter, r *http.Request) error {
	in := &pb.CalculationInput{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Companies.Calculate(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) createApplication(w http.ResponseWriter, r *http.Request) error {
	in := &pb.DraftInput{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Applications.CreateApplication(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 201, v)
}

func page(r *http.Request) (*pb.PageInput, error) {
	q := r.URL.Query()
	offset := int64(0)
	limit := int64(20)
	var e error
	if q.Has("offset") {
		offset, e = strconv.ParseInt(q.Get("offset"), 10, 32)
		if e != nil {
			return nil, platform.Invalid("Некорректный offset")
		}
	}

	if q.Has("limit") {
		limit, e = strconv.ParseInt(q.Get("limit"), 10, 32)
		if e != nil {
			return nil, platform.Invalid("Некорректный limit")
		}
	}

	for _, k := range []string{"fromDate", "toDate"} {
		if q.Get(k) != "" {
			if _, e = time.Parse(time.RFC3339, q.Get(k)); e != nil {
				return nil, platform.Invalid("Дата должна быть RFC3339")
			}
		}
	}

	if offset < 0 || limit < 1 || limit > 100 {
		return nil, platform.Invalid("offset >= 0; limit от 1 до 100")
	}

	return &pb.PageInput{Offset: int32(offset), Limit: int32(limit), Status: q.Get("status"), Search: q.Get("search"), CompanyId: q.Get("companyId"), Recommendation: q.Get("recommendation"), FromDate: q.Get("fromDate"), ToDate: q.Get("toDate")}, nil
}

func (s *Server) listApplications(w http.ResponseWriter, r *http.Request) error {
	in, e := page(r)
	if e != nil {
		return e
	}

	v, e := s.Applications.ListApplications(r.Context(), in)
	if e != nil {
		return e
	}

	for _, a := range v.Items {
		a.OwnerId = ""
	}

	return reply(w, 200, v)
}

func (s *Server) bankApplications(w http.ResponseWriter, r *http.Request) error {
	in, e := page(r)
	if e != nil {
		return e
	}

	v, e := s.Applications.BankApplications(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) getApplication(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.GetApplication(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) patchApplication(w http.ResponseWriter, r *http.Request) error {
	in := &pb.DraftPatch{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Applications.PatchApplication(r.Context(), &pb.PatchInput{Id: id(r).Id, Patch: in})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) validation(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.Validate(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) submit(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.Submit(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 202, v)
}

func (s *Server) cancel(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.Cancel(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) decision(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.GetDecision(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) accept(w http.ResponseWriter, r *http.Request) error {
	in := &pb.AcceptInput{}
	if e := read(r, in); e != nil {
		return e
	}

	in.Id = id(r).Id
	v, e := s.Applications.Accept(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) decline(w http.ResponseWriter, r *http.Request) error {
	in := &pb.DeclineInput{}
	if r.ContentLength != 0 {
		if e := read(r, in); e != nil {
			return e
		}
	}

	in.Id = id(r).Id
	v, e := s.Applications.Decline(r.Context(), in)
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) issue(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.Issue(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 201, v)
}

func (s *Server) history(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.History(r.Context(), id(r))
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) information(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.InformationRequests(r.Context(), id(r))
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) documents(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.ListDocuments(r.Context(), id(r))
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) scoring(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.GetScoring(r.Context(), id(r))
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) overview(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.GetOverview(r.Context(), id(r))
	if e != nil {
		return e
	}

	b, e := platform.ProtoJSON(v)
	if e != nil {
		return e
	}

	var m map[string]any
	if e = json.Unmarshal(b, &m); e != nil {
		return e
	}

	if c, ok := m["companySnapshot"].(map[string]any); ok {
		flatten(c)
	}

	redact(m)
	write(w, 200, m)
	return nil
}

func (s *Server) decide(w http.ResponseWriter, r *http.Request) error {
	in := &pb.ManualDecision{}
	if e := read(r, in); e != nil {
		return e
	}

	v, e := s.Applications.Decide(r.Context(), &pb.ManualInput{Id: id(r).Id, Decision: in})
	if e != nil {
		return e
	}

	return reply(w, 200, v)
}

func (s *Server) feed(w http.ResponseWriter, r *http.Request) error {
	v, e := s.Applications.Feed(r.Context(), &pb.Empty{})
	if e != nil {
		return e
	}

	return array(w, v)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 11*1024*1024)
	if e := r.ParseMultipartForm(128 * 1024); e != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(e, &tooLarge) {
			return platform.E("FILE_TOO_LARGE", "Превышен размер загрузки", 413)
		}

		return platform.E("INVALID_MULTIPART", "Нужен multipart/form-data", 400)
	}

	defer r.MultipartForm.RemoveAll()
	f, h, e := r.FormFile("file")
	if e != nil {
		return platform.Invalid("Нужен файл в поле file")
	}

	defer f.Close()
	if h.Size > 10485760 {
		return platform.E("FILE_TOO_LARGE", "Файл больше 10 МиБ", 413)
	}

	stream, e := s.Applications.UploadDocument(r.Context())
	if e != nil {
		return e
	}

	m := &pb.UploadMetadata{ApplicationId: id(r).Id, Type: r.FormValue("type"), FileName: h.Filename, ContentType: h.Header.Get("Content-Type"), RequestItemId: r.FormValue("requestItemId")}
	if e = stream.Send(&pb.UploadChunk{Value: &pb.UploadChunk_Metadata{Metadata: m}}); e != nil {
		return e
	}

	buf := make([]byte, 64*1024)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			if e = stream.Send(&pb.UploadChunk{Value: &pb.UploadChunk_Data{Data: buf[:n]}}); e != nil {
				_, recvErr := stream.CloseAndRecv()
				if recvErr != nil {
					return recvErr
				}

				return e
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	v, e := stream.CloseAndRecv()
	if e != nil {
		return e
	}

	return reply(w, 201, v)
}

func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request) error {
	_, e := s.Applications.DeleteDocument(r.Context(), &pb.DocumentID{ApplicationId: id(r).Id, DocumentId: r.PathValue("documentId")})
	if e != nil {
		return e
	}

	w.WriteHeader(204)
	return nil
}

func (s *Server) download(w http.ResponseWriter, r *http.Request) error {
	stream, e := s.Applications.DownloadDocument(r.Context(), &pb.DocumentID{ApplicationId: id(r).Id, DocumentId: r.PathValue("documentId")})
	if e != nil {
		return e
	}

	first, e := stream.Recv()
	if e != nil {
		return e
	}

	w.Header().Set("Content-Type", first.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": first.FileName}))
	if len(first.Data) > 0 {
		if _, e = w.Write(first.Data); e != nil {
			return nil
		}
	}

	for {
		v, e := stream.Recv()
		if e == io.EOF {
			return nil
		}

		if e != nil {
			slog.Error("download interrupted", "error", e)
			return nil
		}

		if _, e = w.Write(v.Data); e != nil {
			return nil
		}
	}
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if e := s.DB.Ping(ctx); e != nil {
		return platform.E("NOT_READY", "База данных недоступна", 503)
	}

	for _, c := range []*grpc.ClientConn{s.CompanyConn, s.ApplicationConn} {
		if _, e := healthpb.NewHealthClient(c).Check(ctx, &healthpb.HealthCheckRequest{}); e != nil {
			return platform.E("NOT_READY", "Сервис недоступен", 503)
		}
	}

	write(w, 200, map[string]string{"status": "ready"})
	return nil
}

func redact(v any) {
	switch x := v.(type) {
	case map[string]any:
		delete(x, "ownerId")
		for _, y := range x {
			redact(y)
		}
	case []any:
		for _, y := range x {
			redact(y)
		}
	}
}

func readAuthJSON(w http.ResponseWriter, r *http.Request, out any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
	decoder.DisallowUnknownFields()
	if e := decoder.Decode(out); e != nil {
		return authJSONError(e)
	}

	var extra any
	if e := decoder.Decode(&extra); e != io.EOF {
		if e != nil {
			return authJSONError(e)
		}

		return platform.E("INVALID_JSON", "Ожидается один JSON-объект", 400)
	}

	return nil
}

func authJSONError(e error) error {
	var tooLarge *http.MaxBytesError
	if errors.As(e, &tooLarge) {
		return platform.E("PAYLOAD_TOO_LARGE", "Тело запроса больше 8 КиБ", 413)
	}

	return platform.E("INVALID_JSON", "Некорректный JSON или неизвестные поля", 400)
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) error {
	if e := s.allowAuth(r, "register", 20); e != nil {
		return e
	}

	var in usecase.Registration
	if e := readAuthJSON(w, r, &in); e != nil {
		return e
	}

	session, e := s.Auth.Register(r.Context(), in)
	if e != nil {
		return e
	}

	write(w, 201, session)
	return nil
}
