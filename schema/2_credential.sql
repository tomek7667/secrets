-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS credential (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    token TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE credential;
-- +goose StatementEnd
