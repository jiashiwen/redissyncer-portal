package main

import (
	"etcdexample/core"
	"etcdexample/global"
	"etcdexample/httpserver"
	"etcdexample/httpserver/router"
	"etcdexample/inspection"
	"etcdexample/node"
	"sync"
)

func main() {

	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()
	//etcdClient := global.GetEtcdClient()
	//defer etcdClient.Close()
	wg := &sync.WaitGroup{}

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
	//wg.Add(1)

	//启动http server
	r := router.RootRouter()
	addr := "0.0.0.0:" + global.RSPViper.GetString("http.port")
	httpserver.StartServer(addr, r)
}
