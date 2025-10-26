package http

import (
	"net/http"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	jwtpkg "github.com/sayeed1999/ride-sharing-golang-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	RegisterUC *usecase.RegisterUsecase
	LoginUC    *usecase.LoginUsecase
	JWTService *jwtpkg.Service
}

func NewAuthHandler(registerUC *usecase.RegisterUsecase, loginUC *usecase.LoginUsecase, jwtService *jwtpkg.Service) *AuthHandler {
	return &AuthHandler{registerUC, loginUC, jwtService}
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

	// generate JWT via injected JWT service
	if h.JWTService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt service not configured"})
		return
	}

	token, err := h.JWTService.GenerateToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
