package main

import (
	"1/internal/delivery"
	"1/internal/delivery/middlewares"
	"1/internal/storage"
	"1/internal/usecase"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := storage.NewUsersStorage(db)
	walletRepo := storage.NewWalletStorage(db)
	transactionRepo := storage.NewTransactionStorage(db)
	inventoryRepo := storage.NewInventoryStorage(db)

	auth := usecase.NewAuthService("secret")
	authUsecase := usecase.NewAuthUseCase(userRepo, walletRepo, auth)
	transactionUseCase := usecase.NewTransactionUsecase(walletRepo, transactionRepo)
	inventoryUseCase := usecase.NewInventoryUsecase(inventoryRepo, walletRepo)

	authHandler := delivery.NewAuthHandler(authUsecase)
	transactionHandler := delivery.NewTransactionHandler(transactionUseCase)
	BuyHandler := delivery.NewBuyHandler(inventoryUseCase)

	r := gin.Default()
	r.POST("/auth", authHandler.Auth)
	protected := r.Group("/api")
	protected.Use(middlewares.JWTAuthMiddleware(auth))
	{
		//protected.GET("/info", GetInfo)
		protected.POST("/sendCoin", transactionHandler.SendTokens)
		protected.GET("/buy/:item_id", BuyHandler.BuyItem)
	}

	r.Run(":8080")
}
