-- name: CreateCredential :one
INSERT INTO credential (
    id,
    token,
    expires_at
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetCredential :one
SELECT *
FROM credential
WHERE token = ?;

-- name: DeleteCredential :exec
DELETE FROM credential
WHERE token = ?;

-- name: ListCredentials :many
SELECT *
FROM credential
ORDER BY created_at DESC;

-- name: UpdateCredential :one
UPDATE credential
SET
    expires_at = ?
WHERE id = ?
RETURNING *;
