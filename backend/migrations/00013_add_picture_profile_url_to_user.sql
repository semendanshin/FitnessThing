-- +goose Up
ALTER TABLE users ADD COLUMN picture_profile_url VARCHAR(255);

-- +goose Down
ALTER TABLE users DROP COLUMN picture_profile_url;