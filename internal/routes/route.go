package routes

import (
	"ewallet/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(users controller.UserController, wallet controller.WalletController, ledger controller.LedgerController) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	v1 := r.Group("/api/v1")
	{
		usr := v1.Group("/users")
		{
			usr.POST("", users.Create)
			usr.PATCH("/:id", users.Update)
			usr.GET("/:id", users.GetByID)
			usr.DELETE("/:id", users.Delete)

		}

		wal := v1.Group("/wallets")
		{
			wal.POST("", wallet.Create)
			wal.PUT("/:walletID/suspend", wallet.SuspendWallet)
			wal.GET("/:walletID", wallet.GetWallet)
			wal.POST("/:walletID/topup", wallet.TopUpWallet)
			wal.POST("/:walletID/pay", wallet.PayingWallet)
			wal.POST("/transfer", wallet.TransferWallet)
		}

		led := v1.Group("/ledgers")
		{
			led.GET("/:walletID/history", ledger.GetByWalletID)
		}
	}

	return r
}
