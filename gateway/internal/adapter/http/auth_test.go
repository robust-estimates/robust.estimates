package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"google.golang.org/grpc"

	pb "robust/contracts/credit/v1"
	"robust/gateway/internal/domain"
	"robust/gateway/internal/usecase"
	"robust/platform"
)

type testUsers struct {
	mu    sync.Mutex
	users map[string]domain.User
}

func (u *testUsers) ByLogin(ctx context.Context, login string) (domain.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	v, ok := u.users[login]
	if !ok {
		return v, domain.ErrUserNotFound
	}

	return v, nil
}

func (u *testUsers) Create(ctx context.Context, login string, v domain.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.users[login]; ok {
		return domain.ErrLoginTaken
	}

	u.users[login] = v

	return nil
}

type testCompanies struct{ pb.CompanyServiceClient }

func (testCompanies) ListCompanies(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.CompanyList, error) {
	return &pb.CompanyList{}, nil
}

func authHandler() (http.Handler, *testUsers) {
	secret := "registration-http-test-secret-32-characters"
	users := &testUsers{users: map[string]domain.User{}}
	server := &Server{Secret: secret, Companies: testCompanies{}, Auth: &usecase.Auth{Users: users, Sign: func(u domain.User) (string, error) { return platform.Token(secret, u.ID, u.DisplayName, u.Role) }}}

	return server.Handler(), users
}

func perform(h http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	return res
}

func TestPublicRegistrationReturnsUsableSession(t *testing.T) {
	h, users := authHandler()
	body := `{"login":"New@Example.com","displayName":"Новый клиент","password":"Test-Register-2026!"}`
	response := perform(h, "POST", "/api/v1/auth/register", body, "")
	if response.Code != 201 {
		t.Fatalf("%d %s", response.Code, response.Body.String())
	}

	if strings.Contains(response.Body.String(), "password") || strings.Contains(response.Body.String(), "Hash") {
		t.Fatal("password data exposed")
	}

	var session usecase.Session
	if e := json.Unmarshal(response.Body.Bytes(), &session); e != nil {
		t.Fatal(e)
	}

	if session.Role != "CLIENT" {
		t.Fatal("wrong role")
	}

	me := perform(h, "GET", "/api/v1/me", "", session.AccessToken)
	if me.Code != 200 || !strings.Contains(me.Body.String(), users.users["new@example.com"].ID) {
		t.Fatal("session does not authorize registered user")
	}

	duplicate := perform(h, "POST", "/api/v1/auth/register", body, "")
	if duplicate.Code != 409 {
		t.Fatal("duplicate accepted")
	}

	login := perform(h, "POST", "/api/v1/auth/login", `{"login":"NEW@EXAMPLE.COM","password":"Test-Register-2026!"}`, "")
	if login.Code != 200 {
		t.Fatal("registered account cannot log in")
	}
}

func TestRegistrationRejectsRoleAndUnknownFields(t *testing.T) {
	for _, body := range []string{`{"login":"new@example.com","displayName":"Клиент","password":"Test-Register-2026!","role":"EMPLOYEE"}`, `{"login":"new@example.com","displayName":"Клиент","password":"Test-Register-2026!","isAdmin":true}`, `{"login":"new@example.com","displayName":"Клиент","password":"Test-Register-2026!"} {}`, `{"login":`} {
		h, users := authHandler()
		response := perform(h, "POST", "/api/v1/auth/register", body, "")
		if response.Code != 400 || len(users.users) != 0 {
			t.Fatal("invalid request saved or wrong status")
		}
	}
}

func TestRegistrationRejectsOversizedBody(t *testing.T) {
	h, _ := authHandler()
	body := `{"login":"new@example.com","displayName":"Клиент","password":"Test-Register-2026!","padding":"` + strings.Repeat("a", 9000) + `"}`
	res := perform(h, "POST", "/api/v1/auth/register", body, "")
	if res.Code != 413 && res.Code != 400 {
		t.Fatalf("oversized payload accepted: %d", res.Code)
	}

	res = perform(h, "POST", "/api/v1/auth/register", `{"login":"new@example.com","displayName":"`+strings.Repeat("a", 9000)+`","password":"Test-Register-2026!"}`, "")
	if res.Code != 413 {
		t.Fatalf("expected 413, got %d", res.Code)
	}
}

func TestRegistrationLimiterHasIndependentBucket(t *testing.T) {
	h, _ := authHandler()
	for i := 0; i < 20; i++ {
		if response := perform(h, "POST", "/api/v1/auth/register", "{", ""); response.Code != 400 {
			t.Fatal("rate limited too early")
		}
	}

	if response := perform(h, "POST", "/api/v1/auth/register", "{", ""); response.Code != 429 {
		t.Fatal("registration is not limited")
	}

	if response := perform(h, "POST", "/api/v1/auth/login", `{"login":"missing@example.com","password":"Test-Register-2026!"}`, ""); response.Code != 401 {
		t.Fatal("registration exhausted login limiter")
	}
}
