package middleware

import (
	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		response.Internal(c, c.Errors.Last().Error())
	}
}
