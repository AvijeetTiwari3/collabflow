package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/middleware"
	"github.com/sde/collabflow/internal/model"
	"github.com/sde/collabflow/internal/repository"
	"github.com/sde/collabflow/internal/service"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates an auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user_exists", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration_failed", "message": err.Error()})
		return
	}

	// Set HTTPOnly cookie for refresh token (XSS protection)
	c.SetCookie("refresh_token", resp.Tokens.RefreshToken, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login_failed", "message": err.Error()})
		return
	}

	// Set HTTPOnly cookie for refresh token
	c.SetCookie("refresh_token", resp.Tokens.RefreshToken, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, resp)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Try reading refresh token from HTTPOnly cookie first, then fallback to JSON body
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_refresh_token"})
			return
		}
		refreshToken = body.RefreshToken
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh_token", "message": err.Error()})
		return
	}

	// Rotate HTTPOnly cookie
	c.SetCookie("refresh_token", tokens.RefreshToken, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, tokens)
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if exists {
		if uid, ok := userID.(uuid.UUID); ok {
			_ = h.authService.Logout(c.Request.Context(), uid)
		}
	}

	// Clear cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserIDKey)
	email, _ := c.Get(middleware.ContextUserEmailKey)
	name, _ := c.Get(middleware.ContextUserNameKey)

	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"email": email,
		"name":  name,
	})
}
