package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"redissyncer-portal/logger"
	"strconv"
)

func main() {
	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()
	cli := global.GetEtcdClient()
	defer cli.Close()

	fmt.Println("--------- test EtcdPaginte -------")
	//pageinte, err := resourceutils.NewEtcdPaginte(cli, "key", 20)
	//
	//if err != nil {
	//	panic(err)
	//}

	//for {
	//	if pageinte.LastPage {
	//		return
	//	}
	//	kvs, err := pageinte.Next()
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	for _, item := range kvs {
	//		fmt.Println("key:", string(item.Key), "   value:", string(item.Value))
	//	}
	//
	//}

	kv := cli.KV
	GetMultipleValuesWithPaginationDemo(context.TODO(), kv)

}

func GetMultipleValuesWithPaginationDemo(ctx context.Context, kv clientv3.KV) {
	pagesize := int64(4)

	global.RSPLog.Sugar().Info("*** GetMultipleValuesWithPaginationDemo()")
	// Delete all keys
	kv.Delete(ctx, "key", clientv3.WithPrefix())

	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)

		fmt.Println(k)
		kv.Put(ctx, k, strconv.Itoa(i))
	}

	optstotal := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithCountOnly(),
	}

	grtotal, _ := kv.Get(ctx, "key", optstotal...)
	total := grtotal.Count
	pages := total / pagesize
	remainder := total % pagesize

	global.RSPLog.Sugar().Info("Total count: ", total, "Pages: ", pages, "Reaminder: ", remainder)

	if pages == 0 && remainder == 0 {
		return
	}

	if pages == 0 {
		opts := []clientv3.OpOption{
			clientv3.WithPrefix(),
		}
		gr, _ := kv.Get(ctx, "key", opts...)

		global.RSPLog.Sugar().Infof("----%d----", 1)
		for _, item := range gr.Kvs {
			fmt.Println(string(item.Key), string(item.Value))
		}
		return
	}

	optsfirst := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(pagesize),
	}

	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(pagesize + 1),
		clientv3.WithFromKey(),
	}

	lastKey := ""
	for i := int64(0); i < pages; i++ {
		if i == int64(0) {
			logger.Logger().Sugar().Infof("----%d----", i+1)
			gr, _ := kv.Get(ctx, "key", optsfirst...)
			lastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
			for _, item := range gr.Kvs {
				fmt.Println(string(item.Key), string(item.Value))
			}
			continue
		}

		//opts = append(opts, clientv3.WithFromKey())
		logger.Logger().Sugar().Infof("----%d----", i+1)
		gr, _ := kv.Get(ctx, lastKey, opts...)
		lastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
		for _, item := range gr.Kvs[1:] {
			fmt.Println(string(item.Key), string(item.Value))
		}
	}

	if remainder > 0 {
		logger.Logger().Sugar().Infof("----%d----", pages+1)
		gr, _ := kv.Get(ctx, lastKey, opts...)
		for _, item := range gr.Kvs[1:] {
			fmt.Println(string(item.Key), string(item.Value))
		}
	}

}
