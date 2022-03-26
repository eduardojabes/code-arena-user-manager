-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS arena_user (
    u_id UUID PRIMARY KEY, 
    u_username VARCHAR(32), 
    u_password VARCHAR(64)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS arena_user;
-- +goose StatementEnd