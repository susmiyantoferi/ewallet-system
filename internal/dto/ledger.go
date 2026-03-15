package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type LedgerResponse struct {
	ID          string          `json:"id"`
	WalletID    string          `json:"wallet_id"`
	ReferenceID string          `json:"reference_id"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string          `json:"currency"`
	Type        string          `json:"type"`
	CreatedAt   time.Time       `json:"created_at"`
}
