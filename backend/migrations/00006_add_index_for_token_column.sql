-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY idx_session_token ON sessions (token);

-- +goose NO TRANSACTION
-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_session_token;
