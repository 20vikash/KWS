CREATE TABLE wgpeer (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    public_key VARCHAR(50),
    ip_address VARCHAR(20)
)
