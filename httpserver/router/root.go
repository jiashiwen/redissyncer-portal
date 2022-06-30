package router

import (
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver/middleware"

	"github.com/gin-gonic/gin"
)

// 初始化总路由

func RootRouter() *gin.Engine {

	// var secrets = gin.H{
	// 	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	// 	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	// 	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
	// }

	var Router = gin.Default()
	// 为用户头像和文件提供静态地址
	//Router.StaticFS(global.RSPConfig.Local.Path, httpserver.Dir(global.GVA_CONFIG.Local.Path))

	// 打开就能玩https了
	// Router.Use(middleware.LoadTls())

	//global.RSPLog.Info("use middleware logger")

	// 跨域
	Router.Use(middleware.Cors())
	global.RSPLog.Info("use middleware cors")
	// Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// global.RSPLog.Info("register swagger handler")

	// 方便统一添加路由组前缀 多服务器上线使用
	PublicGroup := Router.Group("")
	{
		// 注册基础功能路由 不做鉴权
		InitBaseRouter(PublicGroup)
		InitTaskRouter(PublicGroup)
		InitNodeRouter(PublicGroup)
		InitDefaultRouter(PublicGroup)
	}

	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.CheckToken())
	{
		InitTestAuthRouter(PrivateGroup)
	}
	// PrivateGroup.Use(gin.BasicAuth(gin.Accounts{
	// 	"foo":    "bar",
	// 	"austin": "1234",
	// 	"lena":   "hello2",
	// 	"manu":   "4321",
	// }))
	// {
	// 	InitTestAuthRouter(PrivateGroup)
	// }

	//PrivateGroup := Router.Group("")
	//PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//{
	//	router.InitApiRouter(PrivateGroup)                   // 注册功能api路由
	//	router.InitJwtRouter(PrivateGroup)                   // jwt相关路由
	//	router.InitUserRouter(PrivateGroup)                  // 注册用户路由
	//	router.InitMenuRouter(PrivateGroup)                  // 注册menu路由
	//	router.InitEmailRouter(PrivateGroup)                 // 邮件相关路由
	//	router.InitSystemRouter(PrivateGroup)                // system相关路由
	//	router.InitCasbinRouter(PrivateGroup)                // 权限相关路由
	//	router.InitCustomerRouter(PrivateGroup)              // 客户路由
	//	router.InitAutoCodeRouter(PrivateGroup)              // 创建自动化代码
	//	router.InitAuthorityRouter(PrivateGroup)             // 注册角色路由
	//	router.InitSimpleUploaderRouter(PrivateGroup)        // 断点续传（插件版）
	//	router.InitSysDictionaryRouter(PrivateGroup)         // 字典管理
	//	router.InitSysOperationRecordRouter(PrivateGroup)    // 操作记录
	//	router.InitSysDictionaryDetailRouter(PrivateGroup)   // 字典详情管理
	//	router.InitFileUploadAndDownloadRouter(PrivateGroup) // 文件上传下载功能路由
	//	router.InitWorkflowProcessRouter(PrivateGroup)       // 工作流相关接口
	//}
	//global.RSPLog.Info("router register success")
	return Router
}
