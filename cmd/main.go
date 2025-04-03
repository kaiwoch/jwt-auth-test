package main

import (
	"1/internal/delivery"
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
	auth := usecase.NewAuthService("secret")
	authUseCase := usecase.NewAuthUseCase(userRepo, walletRepo, auth)
	authHandler := delivery.NewAuthHandler(authUseCase)

	r := gin.Default()
	r.POST("/auth", authHandler.Auth)

	r.Run(":8080")
}
