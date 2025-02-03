-- +goose Up
ALTER TABLE exercise_instances ADD COLUMN position INT NULL;

-- +goose Down
ALTER TABLE exercise_instances DROP COLUMN position;
