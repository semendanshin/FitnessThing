-- +goose NO TRANSACTION
-- +goose Up
CREATE TABLE llm_settings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    base_prompt VARCHAR(255) NOT NULL DEFAULT '',
    variety_level INT NOT NULL DEFAULT 2 CHECK (
        variety_level >= 1
        AND variety_level <= 3
    ),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX CONCURRENTLY llm_settings_user_id_idx ON llm_settings (user_id);
-- +goose Down
DROP TABLE IF EXISTS llm_settings;
DROP INDEX CONCURRENTLY IF EXISTS llm_settings_user_id_idx;