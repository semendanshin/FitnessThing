-- +goose NO TRANSACTION
-- +goose Up
ALTER TABLE expected_sets DROP COLUMN exercise_id;
ALTER TABLE expected_sets ADD COLUMN exercise_log_id UUID NOT NULL;

ALTER TABLE expected_sets ADD CONSTRAINT expected_sets_exercise_log_id_fkey FOREIGN KEY (exercise_log_id) REFERENCES exercise_logs(id) ON DELETE CASCADE;

DROP INDEX CONCURRENTLY IF EXISTS idx_expected_sets_exercise_id;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_expected_sets_exercise_log_id ON expected_sets (exercise_log_id);

-- +goose Down
ALTER TABLE expected_sets DROP CONSTRAINT expected_sets_exercise_log_id_fkey;
ALTER TABLE expected_sets DROP COLUMN exercise_log_id;
ALTER TABLE expected_sets ADD COLUMN exercise_id UUID NOT NULL;
