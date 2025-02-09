-- +goose Up
ALTER TABLE workouts
ADD COLUMN IF NOT EXISTS reasoning TEXT;
ALTER TABLE workouts
ADD COLUMN IF NOT EXISTS is_ai_generated BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE workouts DROP COLUMN IF EXISTS reasoning;
ALTER TABLE workouts DROP COLUMN IF EXISTS is_ai_generated;