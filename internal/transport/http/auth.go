package http

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	domUser "github.com/rlapenok/exchanger/internal/domain/user"
	"github.com/rlapenok/exchanger/internal/uc/auth"
)

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AuthHandler is the handler for the auth endpoints
type AuthHandler struct {
	loginUC *auth.LoginUseCase
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		loginUC: loginUC,
	}
}

// Login handles the login request
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req loginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := auth.LoginInput{
		Name:     req.Name,
		Password: req.Password,
	}

	out, err := h.loginUC.Login(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domUser.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, domUser.ErrInvalidName),
			errors.Is(err, domUser.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	session := sessions.Default(c)
	session.Set("name", out.Name.Value())
	session.Set("role", out.Role.Value())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("HX-Redirect", "/dashboard")
	c.Status(http.StatusOK)
}

// Logout handles the logout request
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("HX-Redirect", "/login")
	c.Status(http.StatusNoContent)
}
