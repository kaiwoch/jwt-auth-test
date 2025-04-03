package usecase

import (
	"1/internal/storage"
	"errors"
)

type TransactionUsecase struct {
	walletRepo      *storage.WalletPosgresStorage
	transactionRepo *storage.TransactionPostgresStorage
}

func NewTransactionUsecase(walletRepo *storage.WalletPosgresStorage, transactionRepo *storage.TransactionPostgresStorage) *TransactionUsecase {
	return &TransactionUsecase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *TransactionUsecase) TransferCoins(fromUserID, toUserID, amount int) error {
	if fromUserID == toUserID {
		return errors.New("cannot send coins to yourself")
	}

	fromBalance, err := u.walletRepo.GetBalance(fromUserID)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return errors.New("not enough balance")
	}

	toBalance, err := u.walletRepo.GetBalance(toUserID)
	if err != nil {
		return err
	}

	err = u.walletRepo.UpdateBalance(fromUserID, fromBalance-amount)
	if err != nil {
		return err
	}

	err = u.walletRepo.UpdateBalance(toUserID, toBalance+amount)
	if err != nil {
		return err
	}

	err = u.transactionRepo.CreateTransaction(fromUserID, toUserID, amount)
	if err != nil {
		return err
	}

	return nil
}
