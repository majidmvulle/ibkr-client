-- name: CreateSession :one
INSERT INTO sessions (
    account_id,
    session_token_encrypted,
    session_token_hash,
    expires_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetSessionByHash :one
SELECT * FROM sessions
WHERE session_token_hash = $1
AND expires_at > NOW()
LIMIT 1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= NOW();

-- name: DeleteSessionByHash :exec
DELETE FROM sessions
WHERE session_token_hash = $1;
