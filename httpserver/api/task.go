package api

import (
	"io/ioutil"
	"net/http"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver/model/response"
	"redissyncer-portal/httpserver/service"

	"github.com/gin-gonic/gin"
)

func TaskCreate(c *gin.Context) {
	// createjson := model.TaskCreateBody{}
	// str := ""
	// c.BindJSON(&str)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
	}

	global.RSPLog.Sugar().Debug(string(body))

	if err := service.CreateTask(string(body)); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	// c.JSON(http.StatusOK, body)
	c.Data(http.StatusOK, "application/json", body)

}
