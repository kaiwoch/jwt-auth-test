package usecase

import (
	"1/internal/storage"
	"errors"
)

type InventoryUsecase struct {
	inventoryRepo *storage.InventoryPostgresStorage
	walletRepo    *storage.WalletPosgresStorage
}

func NewInventoryUsecase(inventoryRepo *storage.InventoryPostgresStorage, walletRepo *storage.WalletPosgresStorage) *InventoryUsecase {
	return &InventoryUsecase{
		inventoryRepo: inventoryRepo,
		walletRepo:    walletRepo,
	}
}

func (u *InventoryUsecase) BuyItem(userID, itemID int) error {
	balance, err := u.walletRepo.GetBalance(userID)
	if err != nil {
		return err
	}
	price, err := u.inventoryRepo.GetItemPrice(itemID)
	if err != nil {
		return err
	}

	if balance < price {
		return errors.New("not enough balance")
	}

	err = u.walletRepo.UpdateBalance(userID, balance-price)
	if err != nil {
		return err
	}

	err = u.inventoryRepo.UpdateUserInventory(userID, itemID)
	if err != nil {
		return err
	}

	return nil
}
