-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS log (
    id TEXT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    event TEXT NOT NULL,
    msg TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE log;
-- +goose StatementEnd
