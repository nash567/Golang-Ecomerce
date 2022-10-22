-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_order(
    id SERIAL  PRIMARY KEY,
    user_id INTEGER REFERENCES "user"(id),
    order_id INTEGER REFERENCES "order"(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_order
-- +goose StatementEnd
