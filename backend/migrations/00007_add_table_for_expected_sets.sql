-- +goose NO TRANSACTION
-- +goose Up
CREATE TABLE expected_sets (
    id            UUID PRIMARY KEY,
    exercise_id   UUID NOT NULL,
    workout_id    UUID NOT NULL,
    reps          INT,
    weight        FLOAT,
    time          INTERVAL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE NO ACTION,
    FOREIGN KEY (workout_id) REFERENCES workouts (id) ON DELETE CASCADE
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_expected_sets_exercise_id ON expected_sets (exercise_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_expected_sets_workout_id ON expected_sets (workout_id);

-- +goose Down
DROP TABLE IF EXISTS expected_sets;

DROP INDEX CONCURRENTLY IF EXISTS idx_expected_sets_exercise_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_expected_sets_workout_id;
