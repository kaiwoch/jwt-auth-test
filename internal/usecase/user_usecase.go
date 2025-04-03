package usecase

import (
	"1/internal/storage"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userStorage   *storage.UsersPostgresStorage
	walletStorage *storage.WalletPosgresStorage
	authService   *AuthService
}

func NewAuthUseCase(userStorage *storage.UsersPostgresStorage, walletStorage *storage.WalletPosgresStorage, authService *AuthService) *UserUsecase {
	return &UserUsecase{userStorage: userStorage, walletStorage: walletStorage, authService: authService}
}

func (u *UserUsecase) LoginOrRegister(username, password string) (string, error) {
	user, err := u.userStorage.GetUserByUsername(username)
	if err == sql.ErrNoRows {

		if err == sql.ErrNoRows {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			user, err := u.userStorage.CreateUser(username, string(hashedPassword))
			if err != nil {
				return "", err
			}
			user, err = u.userStorage.GetUserByUsername(user.Username)
			if err != nil {
				return "", err
			}
			err = u.walletStorage.CreateWallet(user.Id)
			if err != nil {
				return "", err
			}
			return u.authService.GenerateToken(user.Id, username)
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return u.authService.GenerateToken(user.Id, username)
}
