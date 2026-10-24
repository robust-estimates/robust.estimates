set -eu
psql -X -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -v company_password="$COMPANY_DB_PASSWORD" -v application_password="$APPLICATION_DB_PASSWORD" -v auth_password="$AUTH_DB_PASSWORD" <<'SQL'
SELECT format('CREATE ROLE %I LOGIN', v.role_name)
FROM (VALUES ('company_user'), ('application_user'), ('auth_user')) AS v(role_name)
WHERE NOT EXISTS (SELECT 1 FROM pg_roles AS r WHERE r.rolname = v.role_name)
\gexec
SELECT format('ALTER ROLE %I WITH LOGIN PASSWORD %L', 'company_user', :'company_password')
\gexec
SELECT format('ALTER ROLE %I WITH LOGIN PASSWORD %L', 'application_user', :'application_password')
\gexec
SELECT format('ALTER ROLE %I WITH LOGIN PASSWORD %L', 'auth_user', :'auth_password')
\gexec
SELECT format('CREATE DATABASE %I OWNER %I', v.database_name, v.role_name)
FROM (VALUES ('company_db', 'company_user'), ('application_db', 'application_user'), ('auth_db', 'auth_user')) AS v(database_name, role_name)
WHERE NOT EXISTS (SELECT 1 FROM pg_database AS d WHERE d.datname = v.database_name)
\gexec
REVOKE ALL ON DATABASE company_db FROM PUBLIC;
REVOKE ALL ON DATABASE application_db FROM PUBLIC;
REVOKE ALL ON DATABASE auth_db FROM PUBLIC;
SQL
