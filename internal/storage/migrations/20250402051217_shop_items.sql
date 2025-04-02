-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    price INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shop_items;
-- +goose StatementEnd
