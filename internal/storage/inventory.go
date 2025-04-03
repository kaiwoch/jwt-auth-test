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
