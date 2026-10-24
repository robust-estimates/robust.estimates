package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"robust/gateway/internal/domain"
)

type Users struct{ DB *pgxpool.Pool }

func (r *Users) ByLogin(ctx context.Context, login string) (domain.User, error) {
	var u domain.User
	e := r.DB.QueryRow(ctx, "SELECT id,display_name,role,password_hash FROM users WHERE lower(login)=$1", login).Scan(&u.ID, &u.DisplayName, &u.Role, &u.PasswordHash)
	if errors.Is(e, pgx.ErrNoRows) {
		return u, domain.ErrUserNotFound
	}

	return u, e
}

func (r *Users) Seed(ctx context.Context, password string) error {
	if len(password) < 12 {
		return errors.New("DEMO_PASSWORD must contain at least 12 characters")
	}

	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return e
	}

	for _, u := range []struct{ ID, Login, Name, Role string }{{"demo-client", "client@demo.local", "Демо предприниматель", "CLIENT"}, {"demo-employee", "employee@demo.local", "Демо специалист", "EMPLOYEE"}, {"demo-client-other", "other@demo.local", "Другой предприниматель", "CLIENT"}} {
		if _, e = r.DB.Exec(ctx, "INSERT INTO users(id,login,display_name,role,password_hash) VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING", u.ID, u.Login, u.Name, u.Role, string(hash)); e != nil {
			return e
		}
	}

	return nil
}

func (r *Users) Create(ctx context.Context, login string, user domain.User) error {
	_, e := r.DB.Exec(ctx, "INSERT INTO users(id,login,display_name,role,password_hash) VALUES($1,$2,$3,$4,$5)", user.ID, login, user.DisplayName, user.Role, user.PasswordHash)
	var pgError *pgconn.PgError
	if errors.As(e, &pgError) && pgError.Code == "23505" && (pgError.ConstraintName == "users_login_key" || pgError.ConstraintName == "users_login_lower_idx") {
		return domain.ErrLoginTaken
	}

	return e
}
