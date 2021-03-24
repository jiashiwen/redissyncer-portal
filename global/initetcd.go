package global

import (
	"github.com/coreos/etcd/clientv3"
)

var (
	etcdClient *clientv3.Client
)

// 初始化etcdClient
func InitEtcd() {
	var initerr error
	//初始化etcdclient
	var etcdCfg clientv3.Config
	if err := RSPViper.UnmarshalKey("etcd", &etcdCfg); err != nil {
		panic(err)
	}

	etcdClient, initerr = clientv3.New(etcdCfg)
	if initerr != nil {
		panic(initerr)
	}
}
