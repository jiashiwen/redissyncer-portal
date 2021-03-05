package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redissyncer-portal/httpserver/service"
)

type User struct {
	Name     string `json:"name"`
	Password int64  `json:"password"`
}

func TaskCreate(c *gin.Context) {
	json := User{}
	c.BindJSON(&json)

	service.TaskCreate()
	c.JSON(http.StatusOK, gin.H{
		"name":     json.Name,
		"password": json.Password,
	})
}
