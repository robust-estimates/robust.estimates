package usecase

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"robust/gateway/internal/domain"
	platform "robust/platform/core"
)

type Users interface {
	ByLogin(context.Context, string) (domain.User, error)
	Create(context.Context, string, domain.User) error
}

type Auth struct {
	Users Users
	Sign  func(domain.User) (string, error)
}

type Session struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
	Role        string `json:"role"`
}

func (a *Auth) Login(ctx context.Context, login, password string) (Session, error) {
	u, e := a.Users.ByLogin(ctx, strings.ToLower(strings.TrimSpace(login)))
	if e != nil && !errors.Is(e, domain.ErrUserNotFound) {
		return Session{}, e
	}

	hash := u.PasswordHash
	if e != nil {
		hash = "$2a$10$JyjK8Bt2wvl3QkTjj2lMDuHzey26u5Q1MqsiDpC9bDgY3F.RWgJmC"
	}

	check := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if e != nil || check != nil {
		return Session{}, platform.E("INVALID_CREDENTIALS", "Неверный логин или пароль", 401)
	}

	t, e := a.Sign(u)

	return Session{AccessToken: t, TokenType: "Bearer", ExpiresIn: 3600, Role: u.Role}, e
}

type Registration struct {
	Login       string `json:"login"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

func (a *Auth) Register(ctx context.Context, in Registration) (Session, error) {
	login := strings.ToLower(strings.TrimSpace(in.Login))
	name := strings.TrimSpace(in.DisplayName)
	address, e := mail.ParseAddress(login)
	if e != nil || address.Address != login || len(login) > 254 {
		return Session{}, platform.Invalid("Укажите корректный email в поле login")
	}

	if utf8.RuneCountInString(name) < 2 || utf8.RuneCountInString(name) > 100 {
		return Session{}, platform.Invalid("Имя должно содержать от 2 до 100 символов")
	}

	if utf8.RuneCountInString(in.Password) < 12 || len(in.Password) > 72 || strings.TrimSpace(in.Password) == "" {
		return Session{}, platform.Invalid("Пароль должен содержать минимум 12 символов и занимать не больше 72 байт")
	}

	hash, e := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if e != nil {
		return Session{}, e
	}

	user := domain.User{ID: platform.ID(), DisplayName: name, Role: "CLIENT", PasswordHash: string(hash)}
	token, e := a.Sign(user)
	if e != nil {
		return Session{}, e
	}

	if e = a.Users.Create(ctx, login, user); e != nil {
		if errors.Is(e, domain.ErrLoginTaken) {
			return Session{}, platform.E("LOGIN_TAKEN", "Пользователь с таким email уже зарегистрирован", 409)
		}

		return Session{}, e
	}

	return Session{AccessToken: token, TokenType: "Bearer", ExpiresIn: 3600, Role: user.Role}, nil
}
