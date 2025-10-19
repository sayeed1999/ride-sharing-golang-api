package http

import (
	"net/http"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	RegisterUC *usecase.RegisterUsecase
	LoginUC    *usecase.LoginUsecase
}

func NewAuthHandler(registerUC *usecase.RegisterUsecase, loginUC *usecase.LoginUsecase) *AuthHandler {
	return &AuthHandler{registerUC, loginUC}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.RegisterUC.Register(req.Email, req.Password, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.LoginUC.Login(req.Email, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
