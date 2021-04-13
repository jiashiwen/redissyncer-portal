package main

import (
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver"
	"redissyncer-portal/httpserver/router"
	"redissyncer-portal/inspection"
	"redissyncer-portal/node"
	"redissyncer-portal/resourceutils"
	"sync"
)

func main() {

	//fmt.Println(os.Args)

	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

	wg := &sync.WaitGroup{}

	//start CursorGC
	resourceutils.NewEtcdCursorGC(10, 300000)
	wg.Add(1)
	go resourceutils.StartCursorGC(wg)

	//start node
	node := node.NewNode()
	if err := node.Registry(); err != nil {
		panic(err)
	}
	wg.Add(1)
	go node.Start(wg)

	//start inspect server
	inspector := inspection.NewInspector()
	wg.Add(1)
	go inspector.Start(wg)

	//启动http server
	//wg.Add(1)
	r := router.RootRouter()
	addr := "0.0.0.0:" + global.RSPViper.GetString("node.nodeport")
	httpserver.StartServer(addr, r)

}
