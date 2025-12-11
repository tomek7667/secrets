-- name: CreateCertificate :one
INSERT INTO certificate (
    id,
    name,
    cert_type,
    algorithm,
    key_size,
    pem_data,
    metadata
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetCertificate :one
SELECT *
FROM certificate
WHERE name = ?;

-- name: DeleteCertificate :exec
DELETE FROM certificate
WHERE name = ?;

-- name: ListCertificates :many
SELECT *
FROM certificate
ORDER BY created_at DESC;

-- name: UpdateCertificate :one
UPDATE certificate
SET
    pem_data = ?,
    metadata = ?
WHERE name = ?
RETURNING *;

-- name: GetCertificatesByType :many
SELECT *
FROM certificate
WHERE cert_type = ?
ORDER BY created_at DESC;
