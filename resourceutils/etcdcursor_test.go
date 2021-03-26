package resourceutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"strconv"
	//"strings"
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

func TestMap(t *testing.T) {
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

	lease := clientv3.NewLease(cli)
	gReps, _ := lease.Grant(context.Background(), 10)
	time.Sleep(2 * time.Second)
	gReps2, _ := lease.TimeToLive(context.Background(), gReps.ID)
	fmt.Println("remain ttl:", gReps2.TTL)
	time.Sleep(3 * time.Second)
	gReps3, _ := lease.KeepAliveOnce(context.Background(), gReps.ID)
	fmt.Println("remain ttl:", gReps3.TTL)

	valJson, _ := json.Marshal(global.GetNodeInfo())
	fmt.Println(string(valJson))

	cursorMap := GetCursorQueryMap()
	cursor, _ := NewEtcdCursor(cli, "", 6)
	fmt.Println(cursor)
	//
	(*cursorMap)["aa"] = cursor
	fmt.Println((*cursorMap)["aa"])
	fmt.Println((*cursorMap)["bb"])
	//
	//fmt.Println(strings.Split("/tasks/taskid/abc", "/")[3])
}
