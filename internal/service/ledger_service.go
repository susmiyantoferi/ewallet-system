package service

import (
	"context"
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LedgerService interface {
	GetByWalletID(c context.Context, walletID uuid.UUID) ([]*dto.LedgerResponse, error)
}

type ledgerServiceImpl struct {
	LedgerRepo repository.LedgerRepository
	WalletRepo repository.WalletRepository
	Db         *gorm.DB
	Log        *logrus.Logger
	Validate   *validator.Validate
}

func NewLedgerServiceImpl(ledgerRepo repository.LedgerRepository, walletRepo repository.WalletRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) LedgerService {
	return &ledgerServiceImpl{
		LedgerRepo: ledgerRepo,
		WalletRepo: walletRepo,
		Db:         db,
		Log:        log,
		Validate:   validate,
	}
}

func (l *ledgerServiceImpl) GetByWalletID(c context.Context, walletID uuid.UUID) ([]*dto.LedgerResponse, error) {
	db := l.Db.WithContext(c)

	wallet, err := l.WalletRepo.GetByID(db, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	ledger, err := l.LedgerRepo.GetByWalletID(db, wallet.ID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	var responses []*dto.LedgerResponse
	for _, v := range ledger {
		resp := dto.LedgerResponse{
			ID:          v.ID.String(),
			WalletID:    v.WalletID.String(),
			ReferenceID: v.ReferenceID,
			Amount:      v.Amount,
			Currency:    v.Currency,
			Type:        string(v.Type),
			CreatedAt:   v.CreatedAt,
		}

		responses = append(responses, &resp)
	}

	return responses, nil
}
