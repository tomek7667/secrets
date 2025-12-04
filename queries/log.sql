-- name: CreateLog :one
INSERT INTO log (
    id,
    event,
    msg
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: ListLogs :many
SELECT *
FROM log
ORDER BY created_at DESC;

-- name: DeleteLogs :exec
DELETE FROM log;

