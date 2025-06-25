CREATE TABLE pg_service_user (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pg_user_name VARCHAR(100) NOT NULL UNIQUE,
    pg_user_password VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
)

CREATE TABLE pg_service_db (
    id SERIAL PRIMARY KEY,
    pid INTEGER NOT NULL REFERENCES pg_service_user(id) ON DELETE CASCADE,
    db_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (pid, db_name)
)
