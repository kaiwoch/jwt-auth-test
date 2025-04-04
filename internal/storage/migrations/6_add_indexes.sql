-- +goose Up
-- +goose StatementBegin
-- Индексы для таблицы users
CREATE INDEX idx_wallet_user ON user_wallet(user_id);
CREATE INDEX idx_inventory_user ON inventory(user_id);
CREATE INDEX idx_transactions_from ON transactions(from_user_id);
CREATE INDEX idx_transactions_to ON transactions(to_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_wallet_user;
DROP INDEX IF EXISTS idx_transactions_from;
DROP INDEX IF EXISTS idx_transactions_to;
DROP INDEX IF EXISTS idx_transactions_timestamp;
DROP INDEX IF EXISTS idx_inventory_item;
-- +goose StatementEnd