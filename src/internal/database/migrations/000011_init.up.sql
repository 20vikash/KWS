CREATE TABLE tunnels (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain VARCHAR(100),
    is_custom BOOLEAN,
);
