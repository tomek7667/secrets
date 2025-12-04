-- name: CreateSecret :one
INSERT INTO secret (
    id,
    key,
    value
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetSecret :one
SELECT *
FROM secret
WHERE key = ?;

-- name: DeleteSecret :exec
DELETE FROM secret
WHERE key = ?;

-- name: ListSecrets :many
SELECT *
FROM secret
ORDER BY created_at DESC;

-- name: UpdateSecret :one
UPDATE secret
SET
    value = ?
WHERE key = ?
RETURNING *;
