package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	WalletID    uuid.UUID       `gorm:"type:uuid;notnull;index:wallet_id" json:"wallet_id"`
	Wallet      Wallet          `gorm:"foreignKey:WalletID;" json:"wallet"`
	ReferenceID string          `gorm:"type:varchar(100);notnull;uniqueIndex" json:"reference_id"`
	Amount      decimal.Decimal `gorm:"type:decimal(20,2);notnull" json:"amount"`
	Currency    string          `gorm:"type:varchar(3);notnull" json:"currency"`
	Type        LedgerType      `gorm:"type:varchar(20);notnull" json:"type"`
	CreatedAt   time.Time       `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type LedgerType string

const (
	LedgerTypeTopup       LedgerType = "TOPUP"
	LedgerTypePayment     LedgerType = "PAYMENT"
	LedgerTypeTransferOut LedgerType = "TRANSFER_OUT"
	LedgerTypeTransferIn  LedgerType = "TRANSFER_IN"
)
