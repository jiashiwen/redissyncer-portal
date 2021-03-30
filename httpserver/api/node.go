package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redissyncer-portal/httpserver/model/response"
	"redissyncer-portal/httpserver/service"
)

func NodeListAll(c *gin.Context) {
	all, err := service.NodesStatus()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, all)
}
