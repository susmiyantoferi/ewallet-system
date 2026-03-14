package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;notnull;uniqueIndex:user_currency" json:"user_id"`
	User      User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Balance   decimal.Decimal `gorm:"type:decimal(20,2);notnull" json:"balance"`
	Currency  string          `gorm:"type:varchar(3);notnull;uniqueIndex:user_currency" json:"currency"`
	Status    WalletStatus    `gorm:"type:varchar(20);notnull;default:ACTIVE" json:"status"`
	CreatedAt time.Time       `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type WalletStatus string

const (
	WalletStatusActive  WalletStatus = "ACTIVE"
	WalletStatusSuspended WalletStatus = "SUSPENDED"
)
