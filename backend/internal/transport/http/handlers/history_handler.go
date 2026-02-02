package handlers

import (
	"net/http"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/history"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	list *history.ListHistoryUseCase
}

func NewHistoryHandler(list *history.ListHistoryUseCase) *HistoryHandler {
	return &HistoryHandler{list: list}
}

func (h *HistoryHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	items, err := h.list.Execute(c.Request.Context(), history.ListHistoryInput{UserID: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
