-- name: CreatePermission :one
INSERT INTO permission (
    id,
    credential_id,
    secret_key_pattern
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: ListPermissionsByCredential :many
SELECT *
FROM permission
WHERE credential_id = ?
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
