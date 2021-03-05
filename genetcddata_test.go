package main

import (
	"context"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"redissyncer-portal/node"
	"strconv"
	"testing"
)

func TestGenEtcdData(t *testing.T) {
	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

	etcdClient := global.GetEtcdClient()
	defer etcdClient.Close()

	//生成模拟tasks
	for i := 0; i < 10; i++ {
		for j := 0; j < i; j++ {
			key := node.NodeIncludeTasks + strconv.Itoa(i) + "/" + strconv.Itoa(j)
			_, err := etcdClient.Put(context.Background(), key, "{}")
			if err != nil {
				t.Error(err)
			}
		}
	}

}
