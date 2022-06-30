package api

import "github.com/gin-gonic/gin"

func TestAuth(c *gin.Context) {
	c.String(200, "ok")
}
