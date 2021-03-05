package resourceutils

import (
	"context"
	"redissyncer-portal/global"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strconv"
	"testing"
)

func TestNewEtcdPaginte(t *testing.T) {
	cli := global.GetEtcdClient()

	fmt.Println("--------- test EtcdPaginte -------")
	cli.Delete(context.TODO(), "key", clientv3.WithPrefix())

	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)

		fmt.Println(k)
		cli.Put(context.TODO(), k, strconv.Itoa(i))
	}

	for i := 0; i < 30; i++ {
		pageinte, err := NewEtcdPaginte(cli, "key", int64(i+1))

		if err != nil {
			panic(err)
		}

		for {
			if pageinte.LastPage {
				break
			}
			kvs, err := pageinte.Next()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, item := range kvs {
				fmt.Println("key:", string(item.Key), "   value:", string(item.Value))
			}

		}

	}

}

func TestNewEtcdPaginteWithTraverse(t *testing.T) {

	cli := global.GetEtcdClient()

	fmt.Println("--------- test PaginteWithTravers-------")
	cli.Delete(context.TODO(), "key", clientv3.WithPrefix())

	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)
		//fmt.Println(k)
		cli.Put(context.TODO(), k, strconv.Itoa(i))
	}

	pageinte, err := NewEtcdPaginteWithTraverse(cli, "key", 4)

	if err != nil {
		panic(err)
	}

	for _, v := range pageinte.FirstKeyArray {
		fmt.Println(v)
	}

}

func TestEtcdPaginte_GetPage(t *testing.T) {

	cli := global.GetEtcdClient()

	fmt.Println("--------- test Paginte_GetPage--------")
	cli.Delete(context.TODO(), "key", clientv3.WithPrefix())

	// Insert 20 keys
	for i := 0; i < 20; i++ {
		k := fmt.Sprintf("key_%02d", i)
		//fmt.Println(k)
		cli.Put(context.TODO(), k, strconv.Itoa(i))
	}

	pageinte, err := NewEtcdPaginteWithTraverse(cli, "key", 7)

	if err != nil {
		panic(err)
	}

	for i := int64(0); i < pageinte.Pages; i++ {

		fmt.Printf("-----Get Page %d----- \n", i+1)
		kvs, _ := pageinte.GetPage(i + 1)

		for _, v := range kvs {
			fmt.Println(string(v.Key), "|", string(v.Value))
		}
	}

}

//func TestEtcdPaginte_Previous(t *testing.T) {
//	cli := global.GetEtcdClient()
//
//	fmt.Println("--------- test Previous -------")
//	cli.Delete(context.TODO(), "key", clientv3.WithPrefix())
//
//	// Insert 20 keys
//	for i := 0; i < 20; i++ {
//		k := fmt.Sprintf("key_%02d", i)
//		fmt.Println(k)
//		cli.Put(context.TODO(), k, strconv.Itoa(i))
//	}
//
//	pageinte, err := NewEtcdPaginte(cli, "key", 4)
//
//	if err != nil {
//		panic(err)
//	}
//
//	for {
//		if pageinte.LastPage {
//			break
//		}
//		kvs, err := pageinte.Next()
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		for _, item := range kvs {
//			fmt.Println("key:", string(item.Key), "   value:", string(item.Value))
//		}
//
//	}
//
//	fmt.Println("返回倒数第二页")
//	kvs, err := pageinte.Previous()
//
//	if err != nil {
//		fmt.Println(err)
//	}
//	for _, item := range kvs {
//		fmt.Println("key:", string(item.Key), "   value:", string(item.Value))
//	}
//
//}
