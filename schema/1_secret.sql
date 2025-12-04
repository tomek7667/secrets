-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secret (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secret;
-- +goose StatementEnd
