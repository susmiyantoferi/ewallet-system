package dto

import (
	"ewallet/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateWalletReq struct {
	UserID   string `json:"user_id" validate:"required"`
	Currency string `json:"currency" validate:"required,min=1,max=3"`
}

type TopUpWalletReq struct {
	WalletID    uuid.UUID       `json:"wallet_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	ReferenceID string          `json:"reference_id" validate:"required"`
}

type PayingWalletReq struct {
	WalletID    uuid.UUID       `json:"wallet_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	ReferenceID string          `json:"reference_id" validate:"required"`
}

type TransferWalletReq struct {
	FromWallet  uuid.UUID       `json:"from_wallet" validate:"required"`
	ToWallet    uuid.UUID       `json:"to_wallet" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	ReferenceID string          `json:"reference_id" validate:"required"`
}

type WalletResponse struct {
	ID        string              `json:"id"`
	UserID    *string             `json:"user_id,omitempty"`
	User      *UserResponse       `json:"user,omitempty"`
	Balance   decimal.Decimal     `json:"balance"`
	Currency  string              `json:"currency"`
	Status    entity.WalletStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func WalletToResponse(wal *entity.Wallet, user *entity.User) *WalletResponse {
	return &WalletResponse{
		ID: wal.ID.String(),
		User: &UserResponse{
			ID:        user.ID.String(),
			Name:      user.Name,
			Address:   user.Address,
			CreatedAt: user.CreatedAt,
			UpdatedAt: *user.UpdatedAt,
		},
		Balance:   wal.Balance,
		Currency:  wal.Currency,
		Status:    wal.Status,
		CreatedAt: wal.CreatedAt,
		UpdatedAt: *wal.UpdatedAt,
	}
}
