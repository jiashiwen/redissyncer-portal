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
		TaskRouter.POST("create", api.TaskCreate)
		TaskRouter.POST("stop", api.TaskStop)
		TaskRouter.POST("remove", api.TaskRemove)
		TaskRouter.POST("listbyids", api.TaskListByIDs)
		TaskRouter.POST("listbynames", api.TaskListByNames)
		TaskRouter.POST("listbygroupids", api.TaskListByGroupIDs)
		TaskRouter.POST("listall", api.TaskListAll)
	}

	return TaskRouter
}

func InitNodeRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	NodeRouter := Router.Group("node")
	{
		NodeRouter.POST("listall", api.NodeListAll)
	}

	return NodeRouter
}
