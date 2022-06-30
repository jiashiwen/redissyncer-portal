package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("auth")

		if token == "" {

			err := errors.New("invalid token")

			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
