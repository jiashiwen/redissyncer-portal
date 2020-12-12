package utils

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

var etcdClient *EtcdClient
var once sync.Once

//获取单例Etcd Cient
func GetEtcdClient() *clientv3.Client {
	return etcdClient.Client
}

//获取etcd配置
func GetEtcdConfig() *clientv3.Config {
	return &etcdClient.Config
}

//公共日志

//初始化依赖资源

func init() {
	//初始化etcdclient
	config := clientv3.Config{Endpoints: []string{"114.67.112.67:2379"},
		//Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	}
	var err error
	etcdClient, err = NewEtcdClient(config)
	if err != nil {
		panic(err)
	}

}
