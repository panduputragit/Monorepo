package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/auth"
)

type Handler struct {
	auth *auth.Service
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewHandler(authService *auth.Service) *Handler {
	return &Handler{auth: authService}
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	group := router.Group("/auth")
	group.POST("/login", handler.Login)
	group.GET("/validate", handler.ValidateToken)
	group.POST("/refresh", handler.Refresh)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ValidateToken(c *gin.Context) {
	token := bearerToken(c.GetHeader("Authorization"))
	if token == "" {
		token = c.Query("token")
	}

	resp, err := h.auth.ValidateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
