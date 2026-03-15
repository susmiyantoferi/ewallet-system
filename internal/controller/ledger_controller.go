package controller

import (
	"errors"
	"ewallet/internal/dto"
	"ewallet/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LedgerController interface {
	GetByWalletID(c *gin.Context)
}

type ledgerControllerImpl struct {
	LedgerService service.LedgerService
}

func NewLedgerControllerImpl(ledgerService service.LedgerService) LedgerController {
	return &ledgerControllerImpl{
		LedgerService: ledgerService,
	}
}

func (l *ledgerControllerImpl) GetByWalletID(c *gin.Context) {
	walletID := c.Param("walletID")
	id, err := uuid.Parse(walletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("invalid ID"))
		return
	}

	resp, err := l.LedgerService.GetByWalletID(c.Request.Context(), id)
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
