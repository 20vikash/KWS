CREATE TABLE instance (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    volume_name TEXT NOT NULL UNIQUE,
    container_name TEXT NOT NULL UNIQUE,
    is_running BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
