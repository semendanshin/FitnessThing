-- +goose Up
ALTER TABLE expected_sets DROP COLUMN workout_id;

-- +goose Down
ALTER TABLE expected_sets ADD COLUMN workout_id UUID NOT NULL;
AlTER TABLE expected_sets ADD CONSTRAINT expected_sets_workout_id_fkey FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE;
`