package handlers

import (
	"net/http"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/calls"
	"github.com/gin-gonic/gin"
)

type WebRTCHandler struct {
	initiate  *calls.InitiateCallUseCase
	terminate *calls.TerminateCallUseCase
}

func NewWebRTCHandler(initiate *calls.InitiateCallUseCase, terminate *calls.TerminateCallUseCase) *WebRTCHandler {
	return &WebRTCHandler{
		initiate:  initiate,
		terminate: terminate,
	}
}

func (h *WebRTCHandler) Initiate(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "phone_number is required",
		})
		return
	}

	output, err := h.initiate.Execute(c.Request.Context(), calls.InitiateCallInput{
		UserID:      userID,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if errorMsg == "phone_number is required" || errorMsg == "user_id is required" || errorMsg == "invalid phone number" {
			statusCode = http.StatusBadRequest
		} else if errorMsg == "failed to initiate call" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, gin.H{
			"error":   "call_initiation_failed",
			"message": errorMsg,
		})
		return
	}

	resp := gin.H{
		"call_id":    output.CallID,
		"session_id": output.SessionID,
		"sdp_offer":  output.SDPOffer,
		"status":     output.Status,
		"start_time": output.StartTime.Format("2006-01-02T15:04:05Z07:00"),
	}
	if output.VoiceToken != "" {
		resp["voice_token"] = output.VoiceToken
	}
	c.JSON(http.StatusOK, resp)
}

func (h *WebRTCHandler) Terminate(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	var req struct {
		CallID string `json:"call_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "call_id is required",
		})
		return
	}

	output, err := h.terminate.Execute(c.Request.Context(), calls.TerminateCallInput{
		UserID: userID,
		CallID: req.CallID,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if errorMsg == "call not found" {
			statusCode = http.StatusNotFound
		} else if errorMsg == "unauthorized" {
			statusCode = http.StatusForbidden
		} else if errorMsg == "call_id is required" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "call_termination_failed",
			"message": errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"call_id":  output.CallID,
		"duration": output.Duration,
		"status":   output.Status,
	})
}

