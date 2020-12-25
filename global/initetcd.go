package global

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

var (
	etcdClient *clientv3.Client
)

func InitEtcd() {
	//初始化etcdclient
	var etcdCfg clientv3.Config
	err := RSPViper.UnmarshalKey("etcd", &etcdCfg)

	if err != nil {
		panic(err)
	}

	fmt.Println(etcdCfg)

	etcdClient, err = clientv3.New(etcdCfg)
	if err != nil {
		panic(err)

	}
}
