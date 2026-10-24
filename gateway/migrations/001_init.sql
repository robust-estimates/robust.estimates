CREATE TABLE IF NOT EXISTS users (id text PRIMARY KEY, login text UNIQUE NOT NULL, display_name text NOT NULL, role text NOT NULL CHECK(role IN ('CLIENT','EMPLOYEE')), password_hash text NOT NULL);
