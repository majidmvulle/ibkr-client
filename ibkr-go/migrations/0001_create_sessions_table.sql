-- +goose Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id VARCHAR(255) NOT NULL,
    session_token_encrypted BYTEA NOT NULL,  -- AES-256-GCM encrypted token
    session_token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA-256 hash for lookups
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sessions_account_id ON sessions(account_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_token_hash ON sessions(session_token_hash);

-- +goose Down
DROP TABLE sessions;
