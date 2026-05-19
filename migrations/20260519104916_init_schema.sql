-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    refresh_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_used_at TIMESTAMP DEFAULT NOW(),
    ip TEXT NULL,
    user_agent TEXT NULL
);

CREATE UNIQUE INDEX sessions_refresh_hash_unique
ON sessions(refresh_hash);

CREATE INDEX sessions_user_id_idx
ON sessions(user_id);


-- +goose Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;