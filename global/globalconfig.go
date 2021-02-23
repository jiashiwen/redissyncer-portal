package global

import (
	"redissyncer-portal/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

var (
	RSPViper  *viper.Viper
	RSPLog    *zap.Logger
	RSPConfig config.Server
	once      sync.Once
)

//获取单例Etcd Cient
func GetEtcdClient() *clientv3.Client {
	once.Do(func() {
		InitEtcd()
	})
	return etcdClient
}

//获取node information
func GetNodeInfo() *config.NodeInfo {
	var nodeinfo config.NodeInfo
	err := RSPViper.UnmarshalKey("node", &nodeinfo)
	if err != nil {
		panic(err)
	}
	return &nodeinfo
}
