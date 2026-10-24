package core

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	Subject string
	Role    string
	Name    string
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *Error) Error() string { return e.Message }

func E(code, message string, httpStatus int) error { return &Error{code, message, httpStatus} }

func Invalid(message string) error { return E("VALIDATION_FAILED", message, 422) }

func Conflict(message string) error { return E("INVALID_STATE", message, 409) }

func NotFound() error {
	return E("NOT_FOUND", "Объект не найден или недоступен", 404)
}

func Forbidden() error { return E("FORBIDDEN", "Недостаточно прав", 403) }

func ID() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}

	return hex.EncodeToString(b[:])
}

func Now() string { return time.Now().UTC().Format(time.RFC3339Nano) }

type identityKey struct{}

func WithIdentity(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, identityKey{}, c)
}

func Identity(ctx context.Context) Claims { c, _ := ctx.Value(identityKey{}).(Claims); return c }

func Require(ctx context.Context, roles ...string) error {
	for _, r := range roles {
		if Identity(ctx).Role == r {
			return nil
		}
	}

	return Forbidden()
}

type headers struct{ key, version string }

type headersKey struct{}

func WithHeaders(ctx context.Context, key, version string) context.Context {
	return context.WithValue(ctx, headersKey{}, headers{key, version})
}

func Hash(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)

	return hex.EncodeToString(h[:])
}

type Command struct {
	User      string
	Role      string
	Key       string
	Hash      string
	Operation string
	Expected  int64
}

func GetCommand(ctx context.Context, operation string, payload any, version bool) (Command, error) {
	c := Identity(ctx)
	h, _ := ctx.Value(headersKey{}).(headers)
	k := h.key
	if len(k) < 8 || len(k) > 128 {
		return Command{}, E("INVALID_KEY", "Нужен Idempotency-Key длиной 8–128", 400)
	}

	cmd := Command{User: c.Subject, Role: c.Role, Key: k, Operation: operation, Hash: Hash(payload)}
	if version {
		v, e := strconv.ParseInt(strings.Trim(h.version, "\""), 10, 64)
		if e != nil || v < 1 {
			return cmd, E("INVALID_VERSION", "Нужен If-Match с версией заявки", 400)
		}

		cmd.Expected = v
	}

	return cmd, nil
}

func CheckVersion(c Command, v int64) error {
	if c.Expected != v {
		return E("VERSION_CONFLICT", fmt.Sprintf("Версия изменилась: текущая %d", v), 409)
	}

	return nil
}
