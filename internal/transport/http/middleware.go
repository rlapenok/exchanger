package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequestLogger logs request metadata after the handler chain finishes.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		request := c.Request

		c.Next()

		finishedAt := time.Now()
		status := c.Writer.Status()
		attrs := []slog.Attr{
			slog.String("method", request.Method),
			slog.String("path", request.URL.Path),
			slog.Int("status", status),
			slog.Time("started_at", startedAt),
			slog.Time("finished_at", finishedAt),
			slog.Duration("duration", finishedAt.Sub(startedAt)),
			slog.String("client_ip", c.ClientIP()),
		}

		if request.URL.RawQuery != "" {
			attrs = append(attrs, slog.String("query", request.URL.RawQuery))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.String()))
		}

		if status == http.StatusInternalServerError {
			slog.LogAttrs(request.Context(), slog.LevelError, "request completed", attrs...)
			return
		}

		slog.LogAttrs(request.Context(), slog.LevelInfo, "request completed", attrs...)
	}
}

// Recoverer converts panics into 500 responses and lets RequestLogger log them.
func Recoverer() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				_ = c.Error(fmt.Errorf("panic: %v", recovered))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

// Authenticate checks if the user is authenticated.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		name, okName := session.Get("name").(string)
		role, okRole := session.Get("role").(string)
		if !okName || !okRole || name == "" || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("name", name)
		c.Set("role", role)

		c.Next()
	}
}

// Authorization allows read requests for any authenticated user.
// Mutating requests require admin role.
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") == "admin" {
			c.Next()
			return
		}

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		}
	}
}

// RequireAdmin allows only admin users.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

// RequireOperator allows only operator users.
func RequireOperator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "operator" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
