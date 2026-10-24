package usecase

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"robust/gateway/internal/domain"
	platform "robust/platform/core"
)

type userMemory struct {
	mu          sync.Mutex
	users       map[string]domain.User
	createError error
	readError   error
}

func (r *userMemory) ByLogin(ctx context.Context, login string) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.readError != nil {
		return domain.User{}, r.readError
	}

	u, ok := r.users[login]
	if !ok {
		return u, domain.ErrUserNotFound
	}

	return u, nil
}

func (r *userMemory) Create(ctx context.Context, login string, u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.createError != nil {
		return r.createError
	}

	if _, ok := r.users[login]; ok {
		return domain.ErrLoginTaken
	}

	r.users[login] = u

	return nil
}

func newAuth() (*Auth, *userMemory) {
	repo := &userMemory{users: map[string]domain.User{}}

	return &Auth{Users: repo, Sign: func(u domain.User) (string, error) { return "token-" + u.ID, nil }}, repo
}

func code(e error) int {
	var typed *platform.Error
	if errors.As(e, &typed) {
		return typed.Status
	}

	return 0
}

func validRegistration() Registration {
	return Registration{Login: "  USER@Example.com  ", DisplayName: "  Тестовый клиент  ", Password: "Test-Register-2026!"}
}

func TestRegistrationPersistsHashAndClientRole(t *testing.T) {
	auth, repo := newAuth()
	in := validRegistration()
	s, e := auth.Register(context.Background(), in)
	if e != nil {
		t.Fatal(e)
	}

	u := repo.users["user@example.com"]
	if u.ID == "" || u.Role != "CLIENT" || u.DisplayName != "Тестовый клиент" || s.Role != "CLIENT" || s.AccessToken == "" {
		t.Fatalf("invalid user/session: %+v", u)
	}

	if u.PasswordHash == in.Password || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		t.Fatal("password is not hashed correctly")
	}

	if _, e = auth.Login(context.Background(), " USER@EXAMPLE.COM ", in.Password); e != nil {
		t.Fatal(e)
	}
}

func TestRegistrationValidatesBeforeWriting(t *testing.T) {
	for _, in := range []Registration{{Login: "not-email", DisplayName: "Клиент", Password: "Test-Register-2026!"}, {Login: "a@example.com", DisplayName: " ", Password: "Test-Register-2026!"}, {Login: "a@example.com", DisplayName: "Клиент", Password: "short"}, {Login: "a@example.com", DisplayName: "Клиент", Password: strings.Repeat("я", 37)}, {Login: "Name <a@example.com>", DisplayName: "Клиент", Password: "Test-Register-2026!"}} {
		auth, repo := newAuth()
		if _, e := auth.Register(context.Background(), in); code(e) != 422 {
			t.Fatalf("expected 422: %v", e)
		}

		if len(repo.users) != 0 {
			t.Fatal("invalid account saved")
		}
	}
}

func TestDuplicateEmailIsCaseInsensitive(t *testing.T) {
	auth, repo := newAuth()
	in := validRegistration()
	if _, e := auth.Register(context.Background(), in); e != nil {
		t.Fatal(e)
	}

	in.Login = "USER@example.COM"
	if _, e := auth.Register(context.Background(), in); code(e) != 409 {
		t.Fatalf("expected 409: %v", e)
	}

	if len(repo.users) != 1 {
		t.Fatal("duplicate persisted")
	}
}

func TestConcurrentRegistrationCreatesOneAccount(t *testing.T) {
	auth, repo := newAuth()
	var wg sync.WaitGroup
	out := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := auth.Register(context.Background(), validRegistration()); out <- e }()
	}

	wg.Wait()
	close(out)
	success, conflict := 0, 0
	for e := range out {
		if e == nil {
			success++
		} else if code(e) == 409 {
			conflict++
		} else {
			t.Fatal(e)
		}
	}

	if success != 1 || conflict != 1 || len(repo.users) != 1 {
		t.Fatal("concurrent duplicate registration")
	}
}

func TestSigningFailureDoesNotPersistAccount(t *testing.T) {
	auth, repo := newAuth()
	auth.Sign = func(domain.User) (string, error) { return "", errors.New("signer unavailable") }
	if _, e := auth.Register(context.Background(), validRegistration()); e == nil {
		t.Fatal("expected error")
	}

	if len(repo.users) != 0 {
		t.Fatal("account saved without a session")
	}
}

func TestRegistrationAndLoginExposeNoStorageErrorsAsCredentials(t *testing.T) {
	auth, repo := newAuth()
	failure := errors.New("database unavailable")
	repo.createError = failure
	if _, e := auth.Register(context.Background(), validRegistration()); !errors.Is(e, failure) {
		t.Fatal("storage error swallowed")
	}

	repo.readError = failure
	if _, e := auth.Login(context.Background(), "a@example.com", "Test-Register-2026!"); !errors.Is(e, failure) {
		t.Fatal("login storage failure converted to credentials failure")
	}
}

func TestUnknownAccountReturnsUnauthorized(t *testing.T) {
	auth, _ := newAuth()
	if _, e := auth.Login(context.Background(), "missing@example.com", "Test-Register-2026!"); code(e) != 401 {
		t.Fatal("unknown account accepted")
	}
}
