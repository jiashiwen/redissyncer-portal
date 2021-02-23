package router

import (
	"redissyncer-portal/httpserver/api"
	"github.com/gin-gonic/gin"
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
