-- +goose Up
-- Adds a last_fetched_at column to the feeds table. It should be nullable.
-- +goose StatementBegin
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds DROP COLUMN last_fetched_at;
-- +goose StatementEnd