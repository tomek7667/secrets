-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS certificate (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL UNIQUE,
    cert_type TEXT NOT NULL, -- 'private_key', 'public_key', 'certificate', 'ca_certificate'
    algorithm TEXT NOT NULL, -- 'RSA', 'ECDSA', 'ED25519'
    key_size INTEGER, -- for RSA: 2048, 3072, 4096; for ECDSA: 256, 384, 521
    pem_data TEXT NOT NULL, -- PEM encoded certificate/key data
    metadata TEXT -- JSON metadata (issuer, subject, expiry, etc.)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE certificate;
-- +goose StatementEnd
