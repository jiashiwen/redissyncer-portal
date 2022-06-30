package v1

import (
	"net/http"
	"redissyncer-portal/httpserver/model"
	"redissyncer-portal/httpserver/model/response"
	"redissyncer-portal/httpserver/service"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var loginJson model.Login
	if err := c.ShouldBindJSON(&loginJson); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := service.LoginService(loginJson)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	c.Data(http.StatusOK, "application/json", []byte(resp))

}
