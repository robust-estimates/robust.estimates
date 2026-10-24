CREATE TABLE IF NOT EXISTS companies (id text PRIMARY KEY, owner_id text NOT NULL, data jsonb NOT NULL);
CREATE INDEX IF NOT EXISTS companies_owner_idx ON companies(owner_id);
CREATE TABLE IF NOT EXISTS products (id text PRIMARY KEY, data jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS command_receipts (scope text PRIMARY KEY, payload_hash text NOT NULL, response jsonb NOT NULL, expires_at timestamptz NOT NULL);
