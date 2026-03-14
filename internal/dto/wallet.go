package dto

import (
	"ewallet/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateWalletReq struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Currency string    `json:"currency" validate:"required,min=1,max=3"`
}

type UpdateWalletReq struct {
	WalletID uuid.UUID `json:"wallet_id" validate:"required"`
	Currency string    `json:"currency" validate:"required,min=1,max=3"`
}

type WalletResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	User      UserResponse        `json:"user"`
	Balance   decimal.Decimal     `json:"balance"`
	Currency  string              `json:"currency"`
	Status    entity.WalletStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}
