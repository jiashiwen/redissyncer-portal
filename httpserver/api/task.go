package api

import (
	"io/ioutil"
	"net/http"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver/model"
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
		return
	}

	// c.JSON(http.StatusOK, body)
	c.Data(http.StatusOK, "application/json", []byte(resp))

}

func TaskStart(c *gin.Context) {
	var start model.TaskStart
	if err := c.ShouldBindJSON(&start); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := service.StartTask(start)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	c.Data(http.StatusOK, "application/json", []byte(resp))
}

func TaskStop(c *gin.Context) {
	var stopJSON model.TaskIDBody
	if err := c.ShouldBindJSON(&stopJSON); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := service.StopTaskById(stopJSON.TaskID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	c.Data(http.StatusOK, "application/json", []byte(resp))
}

func TaskRemove(c *gin.Context) {
	var removeJSON model.TaskIDBody
	if err := c.ShouldBindJSON(&removeJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := service.RemoveTask(removeJSON.TaskID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//c.Data(http.StatusOK, "application/json", []byte(response.Ok()))
	response.Ok(c)
}

// TaskListAll 列出集群中的所有任务
func TaskListAll(c *gin.Context) {
	var listAllJSON model.TaskListAll

	if err := c.ShouldBindJSON(&listAllJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	all := service.GetAllTaskStatus(listAllJSON)

	c.JSON(http.StatusOK, all)

}

func TaskListByNodeID(c *gin.Context) {
	var listByNode model.TaskListByNode

	if err := c.ShouldBindJSON(&listByNode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := service.TaskStatusByNodeID(listByNode)

	c.JSON(http.StatusOK, resp)
}

func TaskListByIDs(c *gin.Context) {
	var taskIDsJSON model.TaskListByTaskIDs

	if err := c.ShouldBindJSON(&taskIDsJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var ByName = make(map[string]interface{})
	resp := service.GetTaskStatusByIDs(taskIDsJSON.TaskIDs)
	ByName["result"] = resp
	c.JSON(http.StatusOK, ByName)
	//c.Data(http.StatusOK, "application/json", []byte(resp))
}

func TaskListByNames(c *gin.Context) {
	var namesJSON model.TaskListByTaskNames

	if err := c.ShouldBindJSON(&namesJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ByName = make(map[string]interface{})
	resp := service.GetTaskStatusByName(namesJSON.TaskNames)
	ByName["result"] = resp

	//c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, ByName)

}

func TaskListByGroupIDs(c *gin.Context) {
	var groupIDsJSON model.TaskListByGroupIDs

	if err := c.ShouldBindJSON(&groupIDsJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := service.GetTaskStatusByGroupIDs(groupIDsJSON.GroupIDs)
	c.JSON(http.StatusOK, resp)
}
