package repository

import "ewallet/internal/entity"

type WalletRepository interface {
	Repository[entity.Wallet]
}

type walletRepositoryImpl struct {
	repositoryImpl[entity.Wallet]
}

func NewWalletRepositoryImpl() WalletRepository {
	return &walletRepositoryImpl{}
}
