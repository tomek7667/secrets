-- name: CreatePermission :one
INSERT INTO permission (
    id,
    token_id,
    secret_key_pattern
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetPermission :one
SELECT *
FROM permission
WHERE id = ?;

-- name: ListPermissionsByTokenId :many
SELECT *
FROM permission
WHERE token_id = ?
ORDER BY created_at DESC;

-- name: ListPermissions :many
SELECT *
FROM permission
ORDER BY created_at DESC;

-- name: DeletePermission :exec
DELETE FROM permission
WHERE id = ?;

-- name: UpdatePermission :one
UPDATE permission
SET
    secret_key_pattern = ?
WHERE id = ?
RETURNING *;
