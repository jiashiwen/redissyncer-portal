package router

import (
	"github.com/gin-gonic/gin"
	"redissyncer-portal/httpserver/api"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		//BaseRouter.POST("login", v1.Login)
		//BaseRouter.POST("captcha", v1.Captcha)
		BaseRouter.GET("health", api.Health)

	}

	return BaseRouter
}

func InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	TaskRouter := Router.Group("task")
	{

		TaskRouter.GET("create", api.TaskCreate)
		TaskRouter.POST("create", api.TaskCreate)

	}

	return TaskRouter
}
