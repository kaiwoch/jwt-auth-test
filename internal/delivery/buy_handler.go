package delivery

import (
	"1/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BuyHandler struct {
	inventoryUsecase *usecase.InventoryUsecase
}

func NewBuyHandler(inventoryUsecase *usecase.InventoryUsecase) *BuyHandler {
	return &BuyHandler{inventoryUsecase: inventoryUsecase}
}

func (b *BuyHandler) BuyItem(c *gin.Context) {
	itemIdstring := c.Param("item_id")
	itemId, err := strconv.Atoi(itemIdstring)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request param"})
		return
	}

	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err = b.inventoryUsecase.BuyItem(v.(int), itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "succes"})
}
