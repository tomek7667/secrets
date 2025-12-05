-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS permission (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    token_id TEXT NOT NULL,
    secret_key_pattern TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE permission;
-- +goose StatementEnd
