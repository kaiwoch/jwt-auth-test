-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_wallet (
    id SERIAL PRIMARY KEY,
    user_id INT,
    coin_balance INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_wallet;
-- +goose StatementEnd
