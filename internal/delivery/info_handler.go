package delivery

import (
	"1/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	historyUsecase *usecase.HistoryUsecase
}

func NewInfoHandler(historyUsecase *usecase.HistoryUsecase) *InfoHandler {
	return &InfoHandler{historyUsecase: historyUsecase}
}

func (h *InfoHandler) Info(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := h.historyUsecase.GetInfo(v.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
