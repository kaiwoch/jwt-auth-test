package storage

import (
	"1/internal/storage/entity"
	"database/sql"
)

type UsersPostgresStorage struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersPostgresStorage {
	return &UsersPostgresStorage{db: db}
}

func (u *UsersPostgresStorage) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE username = $1"
	err := u.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersPostgresStorage) CreateUser(username, password string) (*entity.User, error) {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err := u.db.Exec(query, username, password)
	if err != nil {
		return nil, err
	}
	return &entity.User{Username: username, Password: string(password)}, nil
}
