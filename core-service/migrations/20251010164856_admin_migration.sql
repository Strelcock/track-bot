-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
BEGIN;
CREATE TABLE IF NOT EXISTS users_new (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
INSERT INTO users_new (id, name, is_active)
SELECT id, name, is_active FROM users;
DROP TABLE IF EXISTS users;
ALTER TABLE users_new RENAME TO users;
COMMIT;


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE IF EXISTS users DROP COLUMN IF EXISTS is_admin CASCADE;
