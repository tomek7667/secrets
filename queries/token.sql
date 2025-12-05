-- name: CreateToken :one
INSERT INTO token (
    id,
    token,
    expires_at
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetToken :one
SELECT *
FROM token
WHERE token = ?;

-- name: DeleteToken :exec
DELETE FROM token
WHERE token = ?;

-- name: ListTokens :many
SELECT *
FROM token
ORDER BY created_at DESC;

-- name: UpdateToken :one
UPDATE token
SET
    expires_at = ?
WHERE id = ?
RETURNING *;
