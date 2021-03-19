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

	resp, err := service.CreateTask(string(body))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		// return
	}

	// c.JSON(http.StatusOK, body)
	c.Data(http.StatusOK, "application/json", []byte(resp))

}
