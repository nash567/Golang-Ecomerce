-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order"(
    id SERIAL PRIMARY KEY,
    user_id INTEGER  REFERENCES "user"(id),
    product_id INTEGER  REFERENCES product(product_id),
    ordered_at TIMESTAMP NOT NULL,
    price INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    payment_method INTEGER NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order"
-- +goose StatementEnd
