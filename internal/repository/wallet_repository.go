package repository

import (
	"ewallet/internal/entity"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Repository[entity.Wallet]
	GetUserAndCurrency(db *gorm.DB, userID, currency string) (*entity.Wallet, error)
}

type walletRepositoryImpl struct {
	repositoryImpl[entity.Wallet]
}

func NewWalletRepositoryImpl() WalletRepository {
	return &walletRepositoryImpl{}
}

func (w *walletRepositoryImpl) GetUserAndCurrency(db *gorm.DB, userID, currency string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	return &wallet, db.Where("user_id = ? AND currency = ?", userID, currency).Take(&wallet).Error
}