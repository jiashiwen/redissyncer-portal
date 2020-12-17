package main

import (
	"context"
	"etcdexample/global"
	"etcdexample/inspection"
	"fmt"
	"sync"
	"time"
)

func main() {
	cli := global.GetEtcdClient()
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//put key
	resp, _ := cli.Put(ctx, "sample_key", "sample_value")
	fmt.Println(resp)
	//get key
	result, _ := cli.Get(ctx, "sample_key")
	for _, ev := range result.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}

	// watch key:q1mi change
	//rch := cli.Watch(context.Background(), "q1mi") // <-chan WatchResponse
	//for wresp := range rch {
	//	for _, ev := range wresp.Events {
	//		fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	//	}
	//}
	cancel()

	inspector := inspection.NewInspector()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go inspector.Start(wg)
	wg.Wait()

}
