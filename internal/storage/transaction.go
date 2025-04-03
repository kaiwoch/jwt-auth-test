package storage

import (
	"1/internal/storage/entity"
	"database/sql"
)

type TransactionPostgresStorage struct {
	db *sql.DB
}

func NewTransactionStorage(db *sql.DB) *TransactionPostgresStorage {
	return &TransactionPostgresStorage{db: db}
}

func (t *TransactionPostgresStorage) CreateTransaction(fromUserID, toUserID, amount int) error {
	query := "INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)"
	_, err := t.db.Exec(query, fromUserID, toUserID, amount)
	return err
}

func (t *TransactionPostgresStorage) GetTransactionHistory(userID int) (entity.TransactionHistory, error) {
	query := "SELECT users.username, transactions.amount FROM transactions INNER JOIN users ON users.id = transactions.from_user_id WHERE transactions.to_user_id = $1"
	rows, err := t.db.Query(query, userID)
	if err != nil {
		return entity.TransactionHistory{}, err
	}
	defer rows.Close()

	var received []entity.Received
	for rows.Next() {
		var r entity.Received
		err := rows.Scan(&r.FromUser, &r.Amount)
		if err != nil {
			return entity.TransactionHistory{}, err
		}
		received = append(received, r)
	}

	query = "SELECT users.username, transactions.amount FROM transactions JOIN users ON users.id = transactions.to_user_id WHERE transactions.from_user_id = $1"
	rows, err = t.db.Query(query, userID)
	if err != nil {
		return entity.TransactionHistory{}, err
	}
	defer rows.Close()

	var send []entity.Send
	for rows.Next() {
		var s entity.Send
		err := rows.Scan(&s.ToUser, &s.Amount)
		if err != nil {
			return entity.TransactionHistory{}, err
		}
		send = append(send, s)
	}
	return entity.TransactionHistory{Received: received, Send: send}, nil
}
