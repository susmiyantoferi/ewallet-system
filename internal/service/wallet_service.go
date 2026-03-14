package service

import (
	"context"
	"ewallet/internal/dto"

	"github.com/google/uuid"
)

type WalletService interface {
	Create(c context.Context, req *dto.CreateWalletReq) error
	Update(c context.Context, req *dto.CreateWalletReq) error
	Delete(c context.Context, req *dto.CreateWalletReq) error
	GetWallet(c context.Context, userID uuid.UUID) error
}