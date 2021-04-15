package httpserver

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"redissyncer-portal/global"
	"syscall"
	"time"
)

func StartServer(address string, router *gin.Engine) {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20

	s.BeforeBegin = func(add string) {
		pidMap := make(map[string]int)
		// 记录pid
		pid := syscall.Getpid()
		pidMap["pid"] = pid

		pidYaml, _ := yaml.Marshal(pidMap)
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		if err := ioutil.WriteFile(dir+"/pid", pidYaml, 0664); err != nil {
			global.RSPLog.Sugar().Error(err)
		}
		global.RSPLog.Sugar().Infof("Actual pid is %d", pid)
	}
	//s.SignalHooks()
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
