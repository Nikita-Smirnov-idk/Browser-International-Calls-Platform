package handlers

import (
	"net/http"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/calls"
	"github.com/gin-gonic/gin"
)

type CallsHandler struct {
	start *calls.StartCallUseCase
	end   *calls.EndCallUseCase
}

func NewCallsHandler(start *calls.StartCallUseCase, end *calls.EndCallUseCase) *CallsHandler {
	return &CallsHandler{
		start: start,
		end:   end,
	}
}

func (h *CallsHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.start.Execute(c.Request.Context(), calls.StartCallInput{
		UserID:      userID,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           output.CallID,
		"phone_number": req.PhoneNumber,
		"start_time":   output.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		"duration":     0,
		"status":       "initiated",
	})
}

func (h *CallsHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	callID := c.Param("id")
	if callID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "call_id is required"})
		return
	}

	var req struct {
		Duration int    `json:"duration" binding:"required"`
		Status   string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.end.Execute(c.Request.Context(), calls.EndCallInput{
		UserID: userID,
		CallID: callID,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "call not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       callID,
		"duration": req.Duration,
		"status":   req.Status,
	})
}

