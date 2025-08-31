-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS hashed_password TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS hashed_password;