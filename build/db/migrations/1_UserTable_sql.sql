-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id serial PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    phone VARCHAR(10) NOT NULL,
    password VARCHAR NOT NULL,
    created_at Timestamp Default current_timestamp,
    updated_at Timestamp Default current_timestamp
)
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "user"
-- +goose StatementEnd