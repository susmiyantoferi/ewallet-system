package main

import (
	"ewallet/infrastructure/config"
	"ewallet/infrastructure/datastore"
	"ewallet/infrastructure/logger"
	"ewallet/internal/controller"
	"ewallet/internal/repository"
	"ewallet/internal/routes"
	"ewallet/internal/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	config, err := config.NewViper()
	if err != nil {
		panic(err)
	}

	log := logger.NewLogrus(&config.Logger)
	db := datastore.NewDatabase(&config.Postgres)
	datastore.NewRedis(&config.Redis)
	validate := validator.New()

	//user
	userRepo := repository.NewUserRepositoryImpl()
	userService := service.NewUserServiceImpl(userRepo, db, validate, log)
	userController := controller.NewUserControllerImpl(userService)

	//wallet
	walletRepo := repository.NewWalletRepositoryImpl()
	walletService := service.NewWalletServiceImpl(walletRepo, userRepo, db, log, validate)
	walletController := controller.NewWalletControllerImpl(walletService)

	//ledger
	ledgerRepo := repository.NewLedgerRepositoryImpl()
	ledgerService := service.NewLedgerServiceImpl(ledgerRepo, walletRepo, db, log, validate)
	ledgerControler := controller.NewLedgerControllerImpl(ledgerService)

	router := routes.NewRouter(userController, walletController, ledgerControler)

	router.Run(config.App.Port)
}
