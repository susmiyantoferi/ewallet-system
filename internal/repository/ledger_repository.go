package repository

import "ewallet/internal/entity"

type LedgerRepository interface {
	Repository[entity.Ledger]
}

type ledgerRepositoryImpl struct {
	repositoryImpl[entity.Ledger]
}

func NewLedgerRepositoryImpl() LedgerRepository {
	return &ledgerRepositoryImpl{}
}
