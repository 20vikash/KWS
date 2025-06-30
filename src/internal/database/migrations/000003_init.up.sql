CREATE TABLE instance (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ins_user VARCHAR(100),
    ins_password VARCHAR(100),
    volume_name TEXT NOT NULL UNIQUE,
    container_name TEXT NOT NULL UNIQUE,
    instance_type TEXT NOT NULL CHECK (instance_type IN ('core')),
    is_running BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
