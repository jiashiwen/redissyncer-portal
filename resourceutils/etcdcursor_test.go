package resourceutils

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"strconv"
	"testing"
	"time"
)

var config string = "../config.yaml"

func TestEtcdCursor(t *testing.T) {
	global.RSPViper = core.Viper(config)

	var etcdCfg clientv3.Config
	if err := global.RSPViper.UnmarshalKey("etcd", &etcdCfg); err != nil {
		panic(err)
	}
	etcdClient, err := NewEtcdClient(etcdCfg)
	if err != nil {
		t.Error(err)
		return
	}
	cli := etcdClient.Client
	defer cli.Close()

	fmt.Println("--------- test EtcdCursor -------")
	cli.Delete(context.TODO(), "key", clientv3.WithPrefix())

	fmt.Println("--------- Generate keys -------")
	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%03d", i)
		if _, err := cli.Put(context.TODO(), k, strconv.Itoa(i)); err != nil {
			t.Error(err)
		}
		fmt.Println(k)
	}

	cursor, err := NewEtcdCursor(cli, "key", 5)
	if err != nil {
		t.Error(err)
		return
	}

	for !cursor.Finish() {

		fmt.Printf("-------currentpage %d------- \n", cursor.EtcdPaginte.CurrentPage)
		kvs, err := cursor.Next()
		if err != nil {
			t.Error(err)
			return
		}

		for _, v := range kvs {
			fmt.Println("key: ", string(v.Key))
			fmt.Println("value: ", string(v.Value))
		}

		time.Sleep(2 * time.Second)

	}

}
