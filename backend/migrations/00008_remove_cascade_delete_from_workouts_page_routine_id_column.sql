-- +goose NO TRANSACTION
-- +goose Up
ALTER TABLE workouts DROP CONSTRAINT workouts_routine_id_fkey;
ALTER TABLE workouts ADD CONSTRAINT workouts_routine_id_fkey FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE workouts DROP CONSTRAINT workouts_routine_id_fkey;
ALTER TABLE workouts ADD CONSTRAINT workouts_routine_id_fkey FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE;
