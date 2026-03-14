package service

import (
	"context"
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/entity"
	"ewallet/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService interface {
	Create(c context.Context, req *dto.CreateWalletReq) (*dto.WalletResponse, error)
	SuspendWallet(c context.Context, walletID uuid.UUID) (*dto.WalletResponse, error)
	GetWallet(c context.Context, walletID uuid.UUID) (*dto.WalletResponse, error)
	TopUpWallet(c context.Context, req *dto.TopUpWalletReq) error
	PayingWallet(c context.Context, req *dto.PayingWalletReq) error
	TransferWallet(c context.Context, req *dto.TransferWalletReq) error
}

type walletServiceImpl struct {
	WalletRepo repository.WalletRepository
	UserRepo   repository.UserRepository
	Db         *gorm.DB
	Log        *logrus.Logger
	Validate   *validator.Validate
}

func NewWalletServiceImpl(walletRepo repository.WalletRepository, userRepo repository.UserRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) WalletService {
	return &walletServiceImpl{
		WalletRepo: walletRepo,
		UserRepo:   userRepo,
		Db:         db,
		Log:        log,
		Validate:   validate,
	}
}

func (w *walletServiceImpl) Create(c context.Context, req *dto.CreateWalletReq) (*dto.WalletResponse, error) {
	db := w.Db.WithContext(c)

	if err := w.Validate.Struct(req); err != nil {
		w.Log.WithError(err).Error("Validation failed")
		return nil, err
	}

	user, err := w.UserRepo.GetByID(db, req.UserID)
	if err != nil {
		w.Log.WithError(err).Error("failed get user by id")
		return nil, gorm.ErrRecordNotFound
	}

	//cek wallet is already
	_, err = w.WalletRepo.GetUserAndCurrency(db, user.ID.String(), req.Currency)
	if err == nil {
		w.Log.Infof("wallet already exist with currency")
		return nil, gorm.ErrDuplicatedKey
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	wallet := entity.Wallet{
		UserID:   user.ID,
		Currency: req.Currency,
		Status:   entity.WalletStatusActive,
	}

	if err := w.WalletRepo.Create(db, &wallet); err != nil {
		w.Log.WithError(err).Error("failed create new wallet")
		return nil, err
	}

	return dto.WalletToResponse(&wallet, user), nil
}

func (w *walletServiceImpl) SuspendWallet(c context.Context, walletID uuid.UUID) (*dto.WalletResponse, error) {
	db := w.Db.WithContext(c)

	wallet, err := w.WalletRepo.GetByID(db, walletID)
	if err != nil {
		w.Log.WithError(err).Error("failed get wallet by id")
		return nil, gorm.ErrRecordNotFound
	}

	if wallet.Status == entity.WalletStatusSuspended {
		w.Log.Infof("wallet already suspended")
		return nil, errors.New("wallet already suspended")
	}

	wallet.Status = entity.WalletStatusSuspended
	if err := w.WalletRepo.Update(db, wallet, wallet.ID); err != nil {
		w.Log.WithError(err).Error("failed suspended wallet")
		return nil, err
	}

	user, err := w.UserRepo.GetByID(db, wallet.UserID)
	if err != nil {
		w.Log.WithError(err).Error("failed get user by id")
		return nil, gorm.ErrRecordNotFound
	}

	return dto.WalletToResponse(wallet, user), nil
}

func (w *walletServiceImpl) GetWallet(c context.Context, walletID uuid.UUID) (*dto.WalletResponse, error) {
	db := w.Db.WithContext(c)

	wallet, err := w.WalletRepo.GetByID(db, walletID)
	if err != nil {
		w.Log.WithError(err).Error("failed get wallet by id")
		return nil, gorm.ErrRecordNotFound
	}

	user, err := w.UserRepo.GetByID(db, wallet.UserID)
	if err != nil {
		w.Log.WithError(err).Error("failed get user by id")
		return nil, gorm.ErrRecordNotFound
	}

	return dto.WalletToResponse(wallet, user), nil
}

func (w *walletServiceImpl) TopUpWallet(c context.Context, req *dto.TopUpWalletReq) error {

	if err := w.Validate.Struct(req); err != nil {
		w.Log.WithError(err).Error("Validation failed")
		return err
	}

	//rounded amount
	amount := req.Amount.Round(2)
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	if err := w.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		//lock wallet
		var wallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, req.WalletID).Error; err != nil {
			w.Log.WithError(err).Error("topup: failed get wallet")
			return gorm.ErrRecordNotFound
		}

		if wallet.Status != entity.WalletStatusActive {
			w.Log.Infof("wallet suspended")
			return errors.New("wallet suspended")
		}

		//cek ref id for idempotency
		var led entity.Ledger
		err := tx.Where("reference_id = ?", req.ReferenceID).First(&led).Error
		if err == nil {
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		update := map[string]any{
			"balance": gorm.Expr("balance + ? ", amount),
		}

		if err := tx.Model(&wallet).Updates(update).Error; err != nil {
			w.Log.WithError(err).Error("topup: failed update balance wallet")
			return err
		}

		//create ledger
		ledger := entity.Ledger{
			WalletID:    wallet.ID,
			ReferenceID: req.ReferenceID,
			Amount:      amount,
			Currency:    wallet.Currency,
			Type:        entity.LedgerTypeTopup,
		}

		if err := tx.Create(&ledger).Error; err != nil {
			w.Log.WithError(err).Error("topup: failed create ledger")
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *walletServiceImpl) PayingWallet(c context.Context, req *dto.PayingWalletReq) error {
	if err := w.Validate.Struct(req); err != nil {
		w.Log.WithError(err).Error("Validation failed")
		return err
	}

	//rounded amount
	amount := req.Amount.Round(2)
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	if err := w.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		//lock wallet
		var wallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, req.WalletID).Error; err != nil {
			w.Log.WithError(err).Error("pay: failed get wallet")
			return gorm.ErrRecordNotFound
		}

		//cek status wallet
		if wallet.Status != entity.WalletStatusActive {
			w.Log.Infof("pay: wallet suspended")
			return errors.New("wallet suspended")
		}

		//cek balance
		if wallet.Balance.LessThan(amount) {
			w.Log.Infof("pay: insufficient balance")
			return errors.New("insufficient balance")
		}

		//cekk for idempotency
		var led entity.Ledger
		err := tx.Where("reference_id = ?", req.ReferenceID).First(&led).Error
		if err == nil {
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		update := map[string]any{
			"balance": gorm.Expr("balance - ?", amount),
		}

		if err := tx.Model(&wallet).Updates(update).Error; err != nil {
			w.Log.WithError(err).Error("pay: failed update balance wallet")
			return err
		}

		ledger := entity.Ledger{
			WalletID:    wallet.ID,
			ReferenceID: req.ReferenceID,
			Amount:      amount.Neg(),
			Currency:    wallet.Currency,
			Type:        entity.LedgerTypePayment,
		}

		if err := tx.Create(&ledger).Error; err != nil {
			w.Log.WithError(err).Error("pay: failed create ledger")
			return err
		}

		return nil

	}); err != nil {
		return err
	}

	return nil
}

func (w *walletServiceImpl) TransferWallet(c context.Context, req *dto.TransferWalletReq) error {
	if err := w.Validate.Struct(req); err != nil {
		w.Log.WithError(err).Error("Validation failed")
		return err
	}

	//rounded amount
	amount := req.Amount.Round(2)
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	//cekk not same wallet tranfer
	if req.FromWallet == req.ToWallet {
		return errors.New("cannot transfer same wallet")
	}

	if err := w.Db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		firstID := req.FromWallet
		secondID := req.ToWallet

		//set first lock n secc lock
		if req.FromWallet.String() > req.ToWallet.String() {
			firstID = req.ToWallet
			secondID = req.FromWallet
		}

		//lock wallet
		var firstWallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&firstWallet, firstID).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed get from wallet")
			return gorm.ErrRecordNotFound
		}

		var seccondWallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&seccondWallet, secondID).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed get wallet")
			return gorm.ErrRecordNotFound
		}

		//sett receiver n sender
		var fromWallet, toWallet entity.Wallet
		if firstWallet.ID == req.FromWallet {
			fromWallet = firstWallet
			toWallet = seccondWallet
		} else {
			fromWallet = seccondWallet
			toWallet = firstWallet
		}

		//cek status wallet
		if fromWallet.Status != entity.WalletStatusActive || toWallet.Status != entity.WalletStatusActive {
			w.Log.Infof("transfer: wallet suspended")
			fmt.Printf("from wallet: %+v", fromWallet.Status)
			fmt.Printf("to wallet: %+v", toWallet.Status)
			return errors.New("wallet suspended")
		}

		//cekk currency
		if fromWallet.Currency != toWallet.Currency {
			w.Log.Infof("transfer currency invalid")
			return errors.New("transfer currency invalid")
		}

		//cek balance
		if fromWallet.Balance.LessThan(amount) {
			w.Log.Infof("transfer: insufficient balance")
			return errors.New("insufficient balance")
		}

		//cekk for idempotency
		var led entity.Ledger
		err := tx.Where("reference_id = ?", req.ReferenceID).First(&led).Error
		if err == nil {
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		mins := map[string]any{
			"balance": gorm.Expr("balance - ?", amount),
		}

		plus := map[string]any{
			"balance": gorm.Expr("balance + ?", amount),
		}

		if err := tx.Model(&fromWallet).Updates(mins).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed reduce balance")
			return err
		}

		if err := tx.Model(&toWallet).Updates(plus).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed plus balance")
			return err
		}

		minsLedger := entity.Ledger{
			WalletID:    fromWallet.ID,
			ReferenceID: req.ReferenceID,
			Amount:      amount.Neg(),
			Currency:    fromWallet.Currency,
			Type:        entity.LedgerTypeTransferOut,
		}

		plusLedger := entity.Ledger{
			WalletID:    toWallet.ID,
			ReferenceID: req.ReferenceID,
			Amount:      amount,
			Currency:    toWallet.Currency,
			Type:        entity.LedgerTypeTransferIn,
		}

		if err := tx.Create(&minsLedger).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed create ledger minus")
			return err
		}

		if err := tx.Create(&plusLedger).Error; err != nil {
			w.Log.WithError(err).Error("transfer: failed create ledger plus")
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
