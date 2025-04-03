package usecase

import (
	"1/internal/storage"
	"errors"
)

type WalletUsecase struct {
	walletStorage *storage.WalletPosgresStorage
}

func NewWalletStorage(walletStorage *storage.WalletPosgresStorage) *WalletUsecase {
	return &WalletUsecase{walletStorage: walletStorage}
}

func (w *WalletUsecase) GetBalance(userID int) (int, error) {
	return w.walletStorage.GetBalance(userID)
}

func (w *WalletUsecase) UpdateBalance(userID int, amount int) error {
	balance, err := w.walletStorage.GetBalance(userID)
	if err != nil {
		return err
	}

	if amount < 0 && balance < -amount {
		return errors.New("insufficient funds")
	}

	return w.walletStorage.UpdateBalance(userID, balance+amount)
}

func (w *WalletUsecase) CreateWallet(userID int) error {
	return w.walletStorage.CreateWallet(userID)
}
