package handlers

import (
	"net/http"
	"strings"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register   *auth.RegisterUseCase
	login      *auth.LoginUseCase
	logout     *auth.LogoutUseCase
	jwtService JWTService
}

type JWTService interface {
	GenerateToken(userID, email string) (string, error)
}

func NewAuthHandler(register *auth.RegisterUseCase, login *auth.LoginUseCase, logout *auth.LogoutUseCase, jwtService JWTService) *AuthHandler {
	return &AuthHandler{
		register:   register,
		login:      login,
		logout:     logout,
		jwtService: jwtService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}
	
	output, err := h.register.Execute(c.Request.Context(), auth.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "registration_error"
		if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
			errorType = "user_already_exists"
		} else if strings.Contains(err.Error(), "invalid") {
			statusCode = http.StatusBadRequest
			errorType = "validation_error"
		}
		c.JSON(statusCode, gin.H{
			"error":   errorType,
			"message": err.Error(),
		})
		return
	}

	token, err := h.jwtService.GenerateToken(output.UserID, output.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_error",
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":    output.UserID,
			"email": output.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}
	
	output, err := h.login.Execute(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil || output == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": "Invalid email or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": output.AccessToken,
		"user": gin.H{
			"id":    output.UserID,
			"email": output.Email,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userID")
	authHeader := c.GetHeader("Authorization")
	token := ""
	
	if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 {
		token = parts[1]
	}

	_ = h.logout.Execute(c.Request.Context(), auth.LogoutInput{
		UserID: userID,
		Token:  token,
	})
	
	c.Status(http.StatusNoContent)
}
