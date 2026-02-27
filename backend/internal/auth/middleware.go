package auth

import (
	"strings"

	"easy-arbitra/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Try cookie first
		if cookie, err := c.Cookie("auth_token"); err == nil && cookie != "" {
			tokenStr = cookie
		}

		// Fall back to Authorization header
		if tokenStr == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenStr = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if tokenStr == "" {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenStr, secret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

func UserIDFromContext(c *gin.Context) int64 {
	if v, exists := c.Get("user_id"); exists {
		return v.(int64)
	}
	return 0
}
