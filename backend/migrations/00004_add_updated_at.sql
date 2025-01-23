-- +goose Up
-- +goose StatementBegin

ALTER TABLE muscle_groups ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE exercises ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE routines ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE exercise_instances ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE sets ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE workouts ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE exercise_logs ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE set_logs ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(); 
ALTER TABLE sessions ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE muscle_groups DROP COLUMN updated_at;
ALTER TABLE exercises DROP COLUMN updated_at;
ALTER TABLE routines DROP COLUMN updated_at;
ALTER TABLE exercise_instances DROP COLUMN updated_at;
ALTER TABLE sets DROP COLUMN updated_at;
ALTER TABLE workouts DROP COLUMN updated_at;
ALTER TABLE exercise_logs DROP COLUMN updated_at;
ALTER TABLE set_logs DROP COLUMN updated_at; 
ALTER TABLE sessions DROP COLUMN updated_at;

-- +goose StatementEnd
