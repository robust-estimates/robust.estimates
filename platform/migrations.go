package platform

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, db *pgxpool.Pool, dir string) error {
	paths, e := filepath.Glob(filepath.Join(dir, "*.sql"))
	if e != nil {
		return e
	}

	sort.Strings(paths)
	tx, e := db.Begin(ctx)
	if e != nil {
		return e
	}

	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(774301)"); e != nil {
		return e
	}

	if _, e = tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations (name text PRIMARY KEY,checksum text NOT NULL,applied_at timestamptz NOT NULL DEFAULT now())"); e != nil {
		return e
	}

	for _, path := range paths {
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}

		h := sha256.Sum256(b)
		checksum := hex.EncodeToString(h[:])
		name := filepath.Base(path)
		var old string
		e = tx.QueryRow(ctx, "SELECT checksum FROM schema_migrations WHERE name=$1", name).Scan(&old)
		if e == nil {
			if old != checksum {
				return fmt.Errorf("migration %s has changed after application", name)
			}

			continue
		}

		if e != pgx.ErrNoRows {
			return e
		}

		if _, e = tx.Exec(ctx, string(b)); e != nil {
			return fmt.Errorf("migration %s: %w", name, e)
		}

		if _, e = tx.Exec(ctx, "INSERT INTO schema_migrations(name,checksum) VALUES($1,$2)", name, checksum); e != nil {
			return e
		}
	}

	return tx.Commit(ctx)
}
