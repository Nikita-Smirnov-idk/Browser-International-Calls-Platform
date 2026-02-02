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

func (h *CallsHandler) Start(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		CountryCode string `json:"countryCode" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.start.Execute(c.Request.Context(), calls.StartCallInput{
		UserID:      userID,
		CountryCode: req.CountryCode,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil || output == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"callId":    output.CallID,
		"rtcConfig": output.RTCConfig,
	})
}

func (h *CallsHandler) End(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		CallID string `json:"callId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.end.Execute(c.Request.Context(), calls.EndCallInput{
		UserID: userID,
		CallID: req.CallID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "call not found"})
		return
	}
	c.Status(http.StatusOK)
}
