package httpserver

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"redissyncer-portal/global"
	"strconv"
	"syscall"
	"time"
)

func StartServer(address string, router *gin.Engine) {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20

	s.BeforeBegin = func(add string) {
		// 记录pid
		pid := syscall.Getpid()
		if err := ioutil.WriteFile("pid", []byte(strconv.Itoa(pid)), 0664); err != nil {
			global.RSPLog.Sugar().Error(err)
		}
		global.RSPLog.Sugar().Infof("Actual pid is %d", pid)
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
