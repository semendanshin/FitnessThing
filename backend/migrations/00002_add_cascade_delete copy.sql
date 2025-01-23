-- +goose Up
-- +goose StatementBegin

ALTER TABLE exercise_muscle_groups DROP CONSTRAINT exercise_muscle_groups_exercise_id_fkey;
ALTER TABLE exercise_muscle_groups ADD CONSTRAINT exercise_muscle_groups_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE;

ALTER TABLE exercise_muscle_groups DROP CONSTRAINT exercise_muscle_groups_muscle_group_id_fkey;
ALTER TABLE exercise_muscle_groups ADD CONSTRAINT exercise_muscle_groups_muscle_group_id_fkey FOREIGN KEY (muscle_group_id) REFERENCES muscle_groups(id) ON DELETE CASCADE;

ALTER TABLE exercise_instances DROP CONSTRAINT exercise_instances_exercise_id_fkey;
ALTER TABLE exercise_instances ADD CONSTRAINT exercise_instances_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE;

ALTER TABLE exercise_instances DROP CONSTRAINT exercise_instances_routine_id_fkey;
ALTER TABLE exercise_instances ADD CONSTRAINT exercise_instances_routine_id_fkey FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE;

ALTER TABLE exercise_logs DROP CONSTRAINT exercise_logs_workout_id_fkey;
ALTER TABLE exercise_logs ADD CONSTRAINT exercise_logs_workout_id_fkey FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE;

ALTER TABLE exercise_logs DROP CONSTRAINT exercise_logs_exercise_id_fkey;
ALTER TABLE exercise_logs ADD CONSTRAINT exercise_logs_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE;

ALTER TABLE routines DROP CONSTRAINT routines_user_id_fkey;
ALTER TABLE routines ADD CONSTRAINT routines_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE sets DROP CONSTRAINT sets_exercise_instance_id_fkey;
ALTER TABLE sets ADD CONSTRAINT sets_exercise_instance_id_fkey FOREIGN KEY (exercise_instance_id) REFERENCES exercise_instances(id) ON DELETE CASCADE;

ALTER TABLE workouts DROP CONSTRAINT workouts_user_id_fkey;
ALTER TABLE workouts ADD CONSTRAINT workouts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE workouts DROP CONSTRAINT workouts_routine_id_fkey;
ALTER TABLE workouts ADD CONSTRAINT workouts_routine_id_fkey FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE;

ALTER TABLE set_logs DROP CONSTRAINT set_logs_exercise_log_id_fkey;
ALTER TABLE set_logs ADD CONSTRAINT set_logs_exercise_log_id_fkey FOREIGN KEY (exercise_log_id) REFERENCES exercise_logs(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
