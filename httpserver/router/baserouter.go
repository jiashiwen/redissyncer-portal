package router

import (
	"redissyncer-portal/httpserver/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		// BaseRouter.POST("login", v1.Login)
		//BaseRouter.POST("captcha", v1.Captcha)
		// BaseRouter.GET("health", api.Health)

	}

	return BaseRouter
}

func InitTestAuthRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("/")
	{

		BaseRouter.GET("testauth", api.AuthResult)

	}

	return BaseRouter
}

func InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	TaskRouter := Router.Group("api/task")
	{
		TaskRouter.POST("create", api.TaskCreate)
		TaskRouter.POST("start", api.TaskStart)
		TaskRouter.POST("stop", api.TaskStop)
		TaskRouter.POST("remove", api.TaskRemove)
		TaskRouter.POST("listbyids", api.TaskListByIDs)
		TaskRouter.POST("listbynames", api.TaskListByNames)
		TaskRouter.POST("listbygroupids", api.TaskListByGroupIDs)
		TaskRouter.POST("listbynode", api.TaskListByNodeID)
		TaskRouter.POST("listall", api.TaskListAll)
		TaskRouter.POST("lastkeyacross", api.TaskGetLastKeyAcross)
	}

	return TaskRouter
}

func InitNodeRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	NodeRouter := Router.Group("api/node")
	{
		NodeRouter.POST("listall", api.NodeListAll)
		NodeRouter.POST("remove", api.RemoveNode)
	}

	return NodeRouter
}
