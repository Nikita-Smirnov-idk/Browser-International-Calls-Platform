package handlers

import (
	"net/http"
	"strconv"
	"time"

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

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var dateFrom *time.Time
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			dateFrom = &t
		}
	}

	var dateTo *time.Time
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			dateTo = &t
		}
	}

	output, err := h.list.Execute(c.Request.Context(), history.ListHistoryInput{
		UserID:   userID,
		Page:     page,
		Limit:    limit,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, output)
}
