package storage

import (
	"1/internal/storage/entity"
	"database/sql"
)

type InventoryPostgresStorage struct {
	db *sql.DB
}

func NewInventoryStorage(db *sql.DB) *InventoryPostgresStorage {
	return &InventoryPostgresStorage{db: db}
}

func (i *InventoryPostgresStorage) GetUserInventory(userID int) ([]entity.Inventory, error) {
	query := "SELECT shop_items.name, inventory.quantity FROM inventory JOIN shop_items ON shop_items.id = inventory.item_id WHERE user_id = $1"
	rows, err := i.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []entity.Inventory
	for rows.Next() {
		var item entity.Inventory
		err = rows.Scan(&item.ItemName, &item.Quantity)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}
	return inventory, nil
}

func (i *InventoryPostgresStorage) GetItemPrice(itemID int) (int, error) {
	var price int
	query := "select price from shop_items where id = $1"

	err := i.db.QueryRow(query, itemID).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (i *InventoryPostgresStorage) UpdateUserInventory(userID, itemID int) error {
	var quantity int
	query := "select quantity from inventory where user_id = $1 and item_id = $2"

	err := i.db.QueryRow(query, userID, itemID).Scan(&quantity)
	if err == sql.ErrNoRows {
		_, err = i.db.Exec("insert into inventory (user_id, item_id, quantity) values ($1, $2, 1)", userID, itemID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	_, err = i.db.Exec("update inventory set quantity = $1 + 1 where user_id = $2 and item_id = $3", quantity, userID, itemID)
	if err != nil {
		return err
	}
	return nil
}
