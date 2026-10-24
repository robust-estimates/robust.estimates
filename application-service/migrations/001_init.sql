CREATE TABLE IF NOT EXISTS applications (id text PRIMARY KEY, owner_id text NOT NULL, company_id text NOT NULL, status text NOT NULL, version bigint NOT NULL, data jsonb NOT NULL, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE INDEX IF NOT EXISTS applications_owner_idx ON applications(owner_id,created_at DESC);
CREATE INDEX IF NOT EXISTS applications_status_idx ON applications(status,created_at DESC);
CREATE TABLE IF NOT EXISTS loans (id text PRIMARY KEY, application_id text UNIQUE NOT NULL REFERENCES applications(id), data jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS processing_jobs (id text PRIMARY KEY, application_id text NOT NULL REFERENCES applications(id), revision bigint NOT NULL, attempts integer NOT NULL DEFAULT 0, state text NOT NULL DEFAULT 'PENDING', next_run_at timestamptz NOT NULL DEFAULT now(), lease_until timestamptz, UNIQUE(application_id,revision));
CREATE INDEX IF NOT EXISTS jobs_pending_idx ON processing_jobs(state,next_run_at);
CREATE TABLE IF NOT EXISTS document_objects (object_key text PRIMARY KEY, state text NOT NULL DEFAULT 'STAGED', created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS command_receipts (scope text PRIMARY KEY, payload_hash text NOT NULL, response jsonb NOT NULL, expires_at timestamptz NOT NULL);
