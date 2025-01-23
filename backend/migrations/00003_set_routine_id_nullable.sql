-- +goose Up
-- +goose StatementBegin

ALTER TABLE workouts ALTER COLUMN routine_id DROP NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workouts ALTER COLUMN routine_id SET NOT NULL;

-- +goose StatementEnd
