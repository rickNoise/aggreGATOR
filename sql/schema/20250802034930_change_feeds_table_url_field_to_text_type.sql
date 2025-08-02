-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds ALTER COLUMN url TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: This is a potentially destructive operation.
-- If any 'title' values are now longer than 255 characters, they will be truncated.
ALTER TABLE feeds ALTER COLUMN url TYPE VARCHAR(255);
-- +goose StatementEnd