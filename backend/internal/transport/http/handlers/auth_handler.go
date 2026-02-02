package handlers

import (
	"net/http"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/use_cases/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register *auth.RegisterUseCase
	login    *auth.LoginUseCase
	logout   *auth.LogoutUseCase
}

func NewAuthHandler(register *auth.RegisterUseCase, login *auth.LoginUseCase, logout *auth.LogoutUseCase) *AuthHandler {
	return &AuthHandler{
		register: register,
		login:    login,
		logout:   logout,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.register.Execute(c.Request.Context(), auth.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.login.Execute(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil || output == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": output.AccessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	_ = c.GetHeader("Authorization")
	_ = h.logout.Execute(c.Request.Context(), auth.LogoutInput{})
	c.Status(http.StatusNoContent)
}
