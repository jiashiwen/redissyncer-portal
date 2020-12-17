package global

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

var (
	etcdClient *clientv3.Client
	once       sync.Once
)

//获取单例Etcd Cient
func GetEtcdClient() *clientv3.Client {
	once.Do(func() {
		initEtcd()
	})
	return etcdClient
}

//获取etcd配置
//func GetEtcdConfig() *clientv3.Config {
//	return &etcdClient.Config
//}

//公共日志

//初始化依赖资源
func initEtcd() {
	//初始化etcdclient
	config := clientv3.Config{Endpoints: []string{"114.67.112.67:2379"},
		//Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	}
	var err error
	etcdClient, err = clientv3.New(config)
	if err != nil {
		panic(err)
	}
}

//func init() {
//	//初始化etcdclient
//	config := clientv3.Config{Endpoints: []string{"114.67.112.67:2379"},
//		//Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
//		DialTimeout: 5 * time.Second,
//	}
//	var err error
//	etcdClient, err = resourceutils.NewEtcdClient(config)
//	if err != nil {
//		panic(err)
//	}
//
//}
