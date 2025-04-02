-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventory (
    user_id INT,
    item_id INT,
    quantity INT DEFAULT 1,
    PRIMARY KEY (user_id, item_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES shop_items(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventory;
-- +goose StatementEnd
