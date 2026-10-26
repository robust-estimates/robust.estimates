package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ReceiptsDDL = `CREATE TABLE IF NOT EXISTS command_receipts (scope text PRIMARY KEY, payload_hash text NOT NULL, response jsonb NOT NULL, expires_at timestamptz NOT NULL);`

func Transaction(ctx context.Context, pool *pgxpool.Pool, c Command, fn func(pgx.Tx) (any, error), out any) error {
	tx, e := pool.Begin(ctx)
	if e != nil {
		return e
	}

	defer tx.Rollback(ctx)
	scope := c.User + "/" + c.Operation + "/" + c.Key
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1, 0))", scope); e != nil {
		return e
	}

	var hash string
	var b []byte
	e = tx.QueryRow(ctx, "SELECT payload_hash,response FROM command_receipts WHERE scope=$1 AND expires_at>now()", scope).Scan(&hash, &b)
	if e == nil {
		if hash != c.Hash {
			return Conflict("Idempotency-Key уже использован с другим запросом")
		}

		if e = json.Unmarshal(b, out); e != nil {
			return e
		}

		return tx.Commit(ctx)
	}

	if !errors.Is(e, pgx.ErrNoRows) {
		return e
	}

	v, e := fn(tx)
	if e != nil {
		return e
	}

	b, e = json.Marshal(v)
	if e != nil {
		return e
	}

	_, e = tx.Exec(ctx, "INSERT INTO command_receipts(scope,payload_hash,response,expires_at) VALUES($1,$2,$3,$4) ON CONFLICT(scope) DO UPDATE SET payload_hash=excluded.payload_hash,response=excluded.response,expires_at=excluded.expires_at", scope, c.Hash, string(b), time.Now().Add(24*time.Hour))
	if e != nil {
		return e
	}

	if e = tx.Commit(ctx); e != nil {
		return e
	}

	return json.Unmarshal(b, out)
}
