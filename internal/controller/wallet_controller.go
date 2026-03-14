package controller

import (
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletController interface {
	Create(c *gin.Context)
	SuspendWallet(c *gin.Context)
	GetWallet(c *gin.Context)
	TopUpWallet(c *gin.Context)
	PayingWallet(c *gin.Context)
	TransferWallet(c *gin.Context)
}

type walletControllerImpl struct {
	WalletService service.WalletService
}

func NewWalletControllerImpl(walletService service.WalletService) WalletController {
	return &walletControllerImpl{
		WalletService: walletService,
	}
}

func (w *walletControllerImpl) Create(c *gin.Context) {
	var req dto.CreateWalletReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	resp, err := w.WalletService.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, dto.ErrorResponse("wallet already exist with currency"))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(resp))

}

func (w *walletControllerImpl) SuspendWallet(c *gin.Context) {
	walletID := c.Param("walletID")
	id, err := uuid.Parse(walletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	resp, err := w.WalletService.SuspendWallet(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "wallet") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse("wallet already suspended"))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resp))
}

func (w *walletControllerImpl) GetWallet(c *gin.Context) {
	walletID := c.Param("walletID")
	id, err := uuid.Parse(walletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	resp, err := w.WalletService.GetWallet(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resp))
}

func (w *walletControllerImpl) TopUpWallet(c *gin.Context) {
	var req dto.TopUpWalletReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	walletID := c.Param("walletID")
	id, err := uuid.Parse(walletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	req.WalletID = id

	if err := w.WalletService.TopUpWallet(c.Request.Context(), &req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "amount must") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "wallet") {
			c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil))
}

func (w *walletControllerImpl) PayingWallet(c *gin.Context) {
	var req dto.PayingWalletReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	// walletID := c.Param("walletID")
	// id, err := uuid.Parse(walletID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
	// 	return
	// }

	// req.WalletID = id

	if err := w.WalletService.PayingWallet(c.Request.Context(), &req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "amount must") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "wallet") {
			c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil))
}

func (w *walletControllerImpl) TransferWallet(c *gin.Context) {
	var req dto.TransferWalletReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid request body"))
		return
	}

	if err := w.WalletService.TransferWallet(c.Request.Context(), &req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "amount must") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "cannot transfer") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "transfer currency") {
			c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "wallet suspended") {
			c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error()))
			return
		}

		if strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil))
}
