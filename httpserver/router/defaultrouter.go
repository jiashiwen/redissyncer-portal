package router

import (
	"redissyncer-portal/httpserver/api"
	"github.com/gin-gonic/gin"
)

func InitDefaultRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	DefaultRouter := Router.Group("")
	{
		//BaseRouter.POST("login", v1.Login)
		//BaseRouter.POST("captcha", v1.Captcha)
		DefaultRouter.GET("health", api.Health)

	}

	return DefaultRouter
}
