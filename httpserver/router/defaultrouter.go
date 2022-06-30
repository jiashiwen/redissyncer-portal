package router

import (
	"redissyncer-portal/httpserver/api"
	v1 "redissyncer-portal/httpserver/api/v1"

	"github.com/gin-gonic/gin"
)

func InitDefaultRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	DefaultRouter := Router.Group("")
	{
		DefaultRouter.POST("login", v1.Login)
		DefaultRouter.GET("v1/health", v1.Health)
		DefaultRouter.GET("health", api.Health)

	}

	return DefaultRouter
}
