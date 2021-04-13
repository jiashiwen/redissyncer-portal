package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redissyncer-portal/httpserver/model"
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

func RemoveNode(c *gin.Context) {
	var remove model.RemoveNodeModel
	if err := c.ShouldBindJSON(&remove); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := service.RemoveNode(remove); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
