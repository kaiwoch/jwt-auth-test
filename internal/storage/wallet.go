package storage

import (
	"database/sql"
)

type WalletPosgresStorage struct {
	db *sql.DB
}

func NewWalletStorage(db *sql.DB) *WalletPosgresStorage {
	return &WalletPosgresStorage{db: db}
}

func (w *WalletPosgresStorage) GetBalance(userID int) (int, error) {
	var balance int
	query := "SELECT coin_balance FROM user_wallet WHERE user_id = $1"
	err := w.db.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (w *WalletPosgresStorage) UpdateBalance(userID, newBalance int) error {
	query := "UPDATE user_wallet SET coin_balance = $1 WHERE user_id = $2"
	_, err := w.db.Exec(query, newBalance, userID)
	return err
}

func (w *WalletPosgresStorage) CreateWallet(userID int) error {
	query := "INSERT INTO user_wallet (user_id, coin_balance) VALUES ($1, $2)"
	_, err := w.db.Exec(query, userID, 1000)
	return err
}
