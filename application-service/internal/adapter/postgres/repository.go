package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"robust/application-service/internal/domain"
	"robust/platform"
)

type Repository struct{ DB *pgxpool.Pool }

func decode(row pgx.Row) (domain.Aggregate, error) {
	var b []byte
	var a domain.Aggregate
	e := row.Scan(&b)
	if errors.Is(e, pgx.ErrNoRows) {
		return a, platform.NotFound()
	}

	if e != nil {
		return a, e
	}

	e = json.Unmarshal(b, &a)

	return a, e
}

func (r *Repository) Get(ctx context.Context, id string) (domain.Aggregate, error) {
	return decode(r.DB.QueryRow(ctx, "SELECT data FROM applications WHERE id=$1", id))
}

func save(ctx context.Context, tx pgx.Tx, a *domain.Aggregate) error {
	b, e := json.Marshal(a)
	if e != nil {
		return e
	}

	_, e = tx.Exec(ctx, "INSERT INTO applications(id,owner_id,company_id,status,version,data,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT(id) DO UPDATE SET status=excluded.status,version=excluded.version,data=excluded.data,updated_at=excluded.updated_at", a.ID, a.OwnerID, a.CompanyID, a.Status, a.Version, string(b), a.CreatedAt, a.UpdatedAt)

	return e
}

func (r *Repository) Apply(ctx context.Context, cmd platform.Command, id string, fn func(*domain.Aggregate) (domain.Result, error)) (domain.Result, error) {
	var out domain.Result
	e := platform.Transaction(ctx, r.DB, cmd, func(tx pgx.Tx) (any, error) {
		var a domain.Aggregate
		var e error
		if id != "" {
			a, e = decode(tx.QueryRow(ctx, "SELECT data FROM applications WHERE id=$1 FOR UPDATE", id))
			if e != nil {
				return nil, e
			}
		}

		oldStatus := a.Status
		oldDocs := append([]domain.Document{}, a.Documents...)
		result, e := fn(&a)
		if e != nil {
			return nil, e
		}

		if a.ID == "" {
			return nil, errors.New("missing aggregate ID")
		}

		a.Version++
		a.UpdatedAt = platform.Now()
		if a.CreatedAt == "" {
			a.CreatedAt = a.UpdatedAt
		}

		if e = save(ctx, tx, &a); e != nil {
			return nil, e
		}

		if a.Status == "SCORING" && oldStatus != "SCORING" {
			_, e = tx.Exec(ctx, "INSERT INTO processing_jobs(id,application_id,revision) VALUES($1,$2,$3)", platform.ID(), a.ID, a.Version)
			if e != nil {
				return nil, e
			}
		}

		if result.Loan != nil {
			b, _ := json.Marshal(result.Loan)
			if _, e = tx.Exec(ctx, "INSERT INTO loans(id,application_id,data) VALUES($1,$2,$3)", result.Loan.ID, a.ID, string(b)); e != nil {
				return nil, e
			}
		}

		current := map[string]bool{}
		for _, d := range a.Documents {
			current[d.ObjectKey] = true
			if _, e = tx.Exec(ctx, "UPDATE document_objects SET state='ACTIVE' WHERE object_key=$1", d.ObjectKey); e != nil {
				return nil, e
			}
		}

		for _, d := range oldDocs {
			if !current[d.ObjectKey] {
				if _, e = tx.Exec(ctx, "UPDATE document_objects SET state='DELETING' WHERE object_key=$1", d.ObjectKey); e != nil {
					return nil, e
				}
			}
		}

		result.Application = &a.Application

		return result, nil
	}, &out)

	return out, e
}

func (r *Repository) List(ctx context.Context, user string, in domain.PageInput) (domain.Page, error) {
	out := domain.Page{Items: []domain.Application{}, Offset: in.Offset, Limit: in.Limit, Counts: map[string]int64{}}
	where := ` WHERE ($1='' OR owner_id=$1) AND ($2='' OR status=$2) AND ($3='' OR company_id=$3) AND ($4='' OR data->>'number' ILIKE '%'||$4||'%' OR data->>'companyName' ILIKE '%'||$4||'%' OR data->>'inn' ILIKE '%'||$4||'%' OR data->>'productName' ILIKE '%'||$4||'%') AND ($5='' OR data->>'recommendation'=$5) AND ($6='' OR created_at >= NULLIF($6,'')::timestamptz) AND ($7='' OR created_at <= NULLIF($7,'')::timestamptz)`
	args := []any{user, in.Status, in.CompanyID, in.Search, in.Recommendation, in.FromDate, in.ToDate}
	if e := r.DB.QueryRow(ctx, "SELECT count(*) FROM applications"+where, args...).Scan(&out.Total); e != nil {
		return out, e
	}

	args = append(args, in.Limit, in.Offset)
	rows, e := r.DB.Query(ctx, "SELECT data FROM applications"+where+" ORDER BY created_at DESC,id LIMIT $8 OFFSET $9", args...)
	if e != nil {
		return out, e
	}

	defer rows.Close()
	for rows.Next() {
		var b []byte
		var a domain.Aggregate
		if e = rows.Scan(&b); e != nil {
			return out, e
		}

		if e = json.Unmarshal(b, &a); e != nil {
			return out, e
		}

		out.Items = append(out.Items, a.Application)
	}

	if e = rows.Err(); e != nil {
		return out, e
	}

	counts, e := r.DB.Query(ctx, "SELECT status,count(*) FROM applications WHERE ($1='' OR owner_id=$1) AND ($2='' OR company_id=$2) GROUP BY status", user, in.CompanyID)
	if e != nil {
		return out, e
	}

	defer counts.Close()
	for counts.Next() {
		var key string
		var n int64
		if e = counts.Scan(&key, &n); e != nil {
			return out, e
		}

		out.Counts[key] = n
	}

	return out, counts.Err()
}

func (r *Repository) Claim(ctx context.Context) (*domain.Job, error) {
	var j domain.Job
	e := r.DB.QueryRow(ctx, `UPDATE processing_jobs SET state='RUNNING',attempts=attempts+1,lease_until=now()+interval '45 seconds' WHERE id=(SELECT id FROM processing_jobs WHERE (state='PENDING' AND next_run_at<=now()) OR (state='RUNNING' AND lease_until<now()) ORDER BY next_run_at FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING id,application_id,revision,attempts`).Scan(&j.ID, &j.ApplicationID, &j.Revision, &j.Attempts)
	if errors.Is(e, pgx.ErrNoRows) {
		return nil, nil
	}

	return &j, e
}

func (r *Repository) Complete(ctx context.Context, j domain.Job, fn func(*domain.Aggregate) error, retry bool) error {
	tx, e := r.DB.Begin(ctx)
	if e != nil {
		return e
	}

	defer tx.Rollback(ctx)
	var attempts int32
	var state string
	if e = tx.QueryRow(ctx, "SELECT attempts,state FROM processing_jobs WHERE id=$1 FOR UPDATE", j.ID).Scan(&attempts, &state); e != nil {
		return e
	}

	if attempts != j.Attempts || state != "RUNNING" {
		return nil
	}

	a, e := decode(tx.QueryRow(ctx, "SELECT data FROM applications WHERE id=$1 FOR UPDATE", j.ApplicationID))
	if e != nil {
		return e
	}

	if a.Status != "SCORING" || a.Version != j.Revision {
		_, e = tx.Exec(ctx, "UPDATE processing_jobs SET state='DONE' WHERE id=$1", j.ID)
		if e != nil {
			return e
		}

		return tx.Commit(ctx)
	}

	if retry {
		_, e = tx.Exec(ctx, "UPDATE processing_jobs SET state='PENDING',next_run_at=$2,lease_until=NULL WHERE id=$1", j.ID, time.Now().Add(time.Duration(j.Attempts*2)*time.Second))
		if e != nil {
			return e
		}

		return tx.Commit(ctx)
	}

	if e = fn(&a); e != nil {
		return e
	}

	a.Version++
	a.UpdatedAt = platform.Now()
	if e = save(ctx, tx, &a); e != nil {
		return e
	}

	if _, e = tx.Exec(ctx, "UPDATE processing_jobs SET state='DONE' WHERE id=$1", j.ID); e != nil {
		return e
	}

	return tx.Commit(ctx)
}

func (r *Repository) StageObject(ctx context.Context, key string) error {
	_, e := r.DB.Exec(ctx, "INSERT INTO document_objects(object_key) VALUES($1)", key)

	return e
}

func (r *Repository) Garbage(ctx context.Context) ([]string, error) {
	rows, e := r.DB.Query(ctx, "SELECT object_key FROM document_objects WHERE state='DELETING' OR (state='STAGED' AND created_at<now()-interval '1 hour') LIMIT 100")
	if e != nil {
		return nil, e
	}

	defer rows.Close()
	keys := []string{}
	for rows.Next() {
		var key string
		if e = rows.Scan(&key); e != nil {
			return nil, e
		}

		keys = append(keys, key)
	}

	return keys, rows.Err()
}

func (r *Repository) RemoveObject(ctx context.Context, key string) error {
	_, e := r.DB.Exec(ctx, "DELETE FROM document_objects WHERE object_key=$1 AND state!='ACTIVE'", key)

	return e
}
