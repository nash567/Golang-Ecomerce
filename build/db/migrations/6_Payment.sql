-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment(
    id SERIAL  PRIMARY KEY,
    digital BOOLEAN,
    COD BOOLEAN
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment
-- +goose StatementEnd
