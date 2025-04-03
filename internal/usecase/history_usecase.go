package usecase

import (
	"1/internal/storage"
	"1/internal/storage/entity"
)

type HistoryUsecase struct {
	inventoryRepo   *storage.InventoryPostgresStorage
	transactionRepo *storage.TransactionPostgresStorage
	walletRepo      *storage.WalletPosgresStorage
}

func NewHistoryUsecase(inventoryRepo *storage.InventoryPostgresStorage, transactionRepo *storage.TransactionPostgresStorage, walletRepo *storage.WalletPosgresStorage) *HistoryUsecase {
	return &HistoryUsecase{inventoryRepo: inventoryRepo, transactionRepo: transactionRepo, walletRepo: walletRepo}
}

func (h *HistoryUsecase) GetInfo(userID int) (*entity.Info, error) {
	balance, err := h.walletRepo.GetBalance(userID)
	if err != nil {
		return nil, err
	}
	inventory, err := h.inventoryRepo.GetUserInventory(userID)
	if err != nil {
		return nil, err
	}
	history, err := h.transactionRepo.GetTransactionHistory(userID)
	if err != nil {
		return nil, err
	}
	return &entity.Info{Balance: balance, Inventory: inventory, CoinHistory: history}, nil
}
