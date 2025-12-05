-- name: CreateUser :one
INSERT INTO user (
    id,
    username,
    password
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM user
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT *
FROM user
WHERE username = ?;

-- name: DeleteUser :exec
DELETE FROM user
WHERE id = ?;

-- name: ListUsers :many
SELECT *
FROM user
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE user
SET
    password = ?
WHERE id = ?
RETURNING *;
