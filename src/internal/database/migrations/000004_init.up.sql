CREATE TABLE wgpeer (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    public_key VARCHAR(50) NOT NULL UNIQUE,
    ip_address INTEGER NOT NULL UNIQUE
)
