-- list_roles.sql
SELECT rolname FROM pg_roles WHERE rolname != 'pg_signal_backend_pid()' ORDER BY rolname;
