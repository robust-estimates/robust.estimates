package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"robust/company-service/internal/domain"
	"robust/platform"
)

type Repository struct{ DB *pgxpool.Pool }

func (r *Repository) Seed(ctx context.Context) error {
	for _, p := range domain.Products() {
		b, _ := json.Marshal(p)
		if _, e := r.DB.Exec(ctx, "INSERT INTO products(id,data) VALUES($1,$2) ON CONFLICT DO NOTHING", p.ID, string(b)); e != nil {
			return e
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (domain.Company, error) {
	var c domain.Company
	var b []byte
	e := r.DB.QueryRow(ctx, "SELECT data FROM companies WHERE id=$1", id).Scan(&b)
	if errors.Is(e, pgx.ErrNoRows) {
		return c, platform.NotFound()
	}

	if e != nil {
		return c, e
	}

	e = json.Unmarshal(b, &c)

	return c, e
}

func (r *Repository) List(ctx context.Context, user string) ([]domain.Company, error) {
	rows, e := r.DB.Query(ctx, "SELECT data FROM companies WHERE owner_id=$1 ORDER BY id", user)
	if e != nil {
		return nil, e
	}

	defer rows.Close()
	cs := []domain.Company{}
	for rows.Next() {
		var b []byte
		var c domain.Company
		if e = rows.Scan(&b); e != nil {
			return nil, e
		}

		if e = json.Unmarshal(b, &c); e != nil {
			return nil, e
		}

		cs = append(cs, c)
	}

	return cs, rows.Err()
}

func (r *Repository) Products(ctx context.Context) ([]domain.Product, error) {
	rows, e := r.DB.Query(ctx, "SELECT data FROM products ORDER BY id")
	if e != nil {
		return nil, e
	}

	defer rows.Close()
	ps := []domain.Product{}
	for rows.Next() {
		var b []byte
		var p domain.Product
		if e = rows.Scan(&b); e != nil {
			return nil, e
		}

		if e = json.Unmarshal(b, &p); e != nil {
			return nil, e
		}

		ps = append(ps, p)
	}

	return ps, rows.Err()
}

func (r *Repository) Change(ctx context.Context, cmd platform.Command, id string, fn func(*domain.Company) error) (domain.Company, error) {
	var out domain.Company
	e := platform.Transaction(ctx, r.DB, cmd, func(tx pgx.Tx) (any, error) {
		var c domain.Company
		if id != "" {
			var b []byte
			e := tx.QueryRow(ctx, "SELECT data FROM companies WHERE id=$1 FOR UPDATE", id).Scan(&b)
			if errors.Is(e, pgx.ErrNoRows) {
				return nil, platform.NotFound()
			}

			if e != nil {
				return nil, e
			}

			if e = json.Unmarshal(b, &c); e != nil {
				return nil, e
			}
		}

		if e := fn(&c); e != nil {
			return nil, e
		}

		b, e := json.Marshal(c)
		if e != nil {
			return nil, e
		}

		_, e = tx.Exec(ctx, "INSERT INTO companies(id,owner_id,data) VALUES($1,$2,$3) ON CONFLICT(id) DO UPDATE SET data=excluded.data", c.ID, c.OwnerID, string(b))

		return c, e
	}, &out)

	return out, e
}
