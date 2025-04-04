package main

import (
	"1/internal/delivery"
	"1/internal/delivery/middlewares"
	"1/internal/storage"
	"1/internal/usecase"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable&binary_parameters=yes")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)

	userRepo := storage.NewUsersStorage(db)
	walletRepo := storage.NewWalletStorage(db)
	transactionRepo := storage.NewTransactionStorage(db)
	inventoryRepo := storage.NewInventoryStorage(db)

	auth := usecase.NewAuthService("secret")
	authUsecase := usecase.NewAuthUseCase(userRepo, walletRepo, auth)
	transactionUseCase := usecase.NewTransactionUsecase(walletRepo, transactionRepo)
	inventoryUseCase := usecase.NewInventoryUsecase(inventoryRepo, walletRepo)
	historyUsecase := usecase.NewHistoryUsecase(inventoryRepo, transactionRepo, walletRepo)

	authHandler := delivery.NewAuthHandler(authUsecase)
	transactionHandler := delivery.NewTransactionHandler(transactionUseCase)
	BuyHandler := delivery.NewBuyHandler(inventoryUseCase)
	InfoHandler := delivery.NewInfoHandler(historyUsecase)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/auth", authHandler.Auth)
	protected := r.Group("/api")
	protected.Use(middlewares.JWTAuthMiddleware(auth))
	{
		protected.GET("/info", InfoHandler.Info)
		protected.POST("/sendCoin", transactionHandler.SendTokens)
		protected.GET("/buy/:item_id", BuyHandler.BuyItem)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
