package httpserver

import (
	"redissyncer-portal/global"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"syscall"
	"time"
)

func StartServer(address string, router *gin.Engine) {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20

	s.BeforeBegin = func(add string) {
		global.RSPLog.Sugar().Infof("Actual pid is %d", syscall.Getpid())
	}
	global.RSPLog.Error(s.ListenAndServe().Error())
	//return s
}

func initServer(address string, router *gin.Engine) {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20

	s.BeforeBegin = func(add string) {
		global.RSPLog.Sugar().Infof("Actual pid is %d", syscall.Getpid())
	}
	global.RSPLog.Error(s.ListenAndServe().Error())
	//return s
}
