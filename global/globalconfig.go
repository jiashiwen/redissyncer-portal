package global

import (
	"etcdexample/config"
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

//公共日志
