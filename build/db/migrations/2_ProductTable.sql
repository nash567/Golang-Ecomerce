-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product(
    product_id SERIAL  PRIMARY KEY,
    product_name VARCHAR(64) NOT NULL,
    product_price INTEGER NOT NULL,
    product_rating INTEGER NOT NULL
    
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product
-- +goose StatementEnd
