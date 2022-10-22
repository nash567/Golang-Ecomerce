-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS address(
    id SERIAL  PRIMARY KEY,
    user_id INTEGER REFERENCES "user"(id),
    house VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(100) NOT NULL,
    pincode INTEGER NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS address
-- +goose StatementEnd
