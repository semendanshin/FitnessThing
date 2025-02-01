-- +goose Up
ALTER TABLE expected_sets ADD COLUMN set_type VARCHAR(255) NOT NULL;

-- +goose Down
ALTER TABLE expected_sets DROP COLUMN set_type;
