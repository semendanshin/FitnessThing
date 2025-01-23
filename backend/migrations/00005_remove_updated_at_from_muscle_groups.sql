-- +goose Up
-- +goose StatementBegin

ALTER TABLE muscle_groups DROP COLUMN updated_at;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE muscle_groups ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose StatementEnd
