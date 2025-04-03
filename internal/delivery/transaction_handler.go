package delivery

import (
	"1/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase *usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: transactionUsecase}
}

func (h *TransactionHandler) SendTokens(c *gin.Context) {
	var input struct {
		ToUser int `json:"toUser"` // позже надо поменять на username
		Amount int `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.transactionUsecase.TransferCoins(v.(int), input.ToUser, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "succes"})
}
