-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product(
    id SERIAL  PRIMARY KEY,
    product_name VARCHAR(64) NOT NULL,
    price INTEGER NOT NULL,
    Rating INTEGER NOT NULL
    
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product
-- +goose StatementEnd
