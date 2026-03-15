package repository

import (
	"ewallet/internal/entity"

	"gorm.io/gorm"
)

type LedgerRepository interface {
	Repository[entity.Ledger]
	GetByWalletID(db *gorm.DB, walletID string) ([]*entity.Ledger, error)
}

type ledgerRepositoryImpl struct {
	repositoryImpl[entity.Ledger]
}

func NewLedgerRepositoryImpl() LedgerRepository {
	return &ledgerRepositoryImpl{}
}

func (l *ledgerRepositoryImpl) GetByWalletID(db *gorm.DB, walletID string) ([]*entity.Ledger, error) {
	var ledger []*entity.Ledger
	return ledger, db.Where("wallet_id = ? ", walletID).Find(&ledger).Error
}
