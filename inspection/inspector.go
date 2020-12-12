package inspection

import (
	"context"
	"etcdexample/utils"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gofrs/uuid"
	"strconv"

	"sync"
	"time"
)

const (
	LastInspectionTime = "lastinspectiontime"
	InspectLockKey     = "/lock/inspect"
)

type Inspector struct {
	//服务器名称
	ServerName string

	//etcd 客户端
	EtcdClient *clientv3.Client

	//是否执行巡检
	ExecStatus bool

	//轮询Ticker
	InspectTicker *time.Ticker

	//巡检状态变更锁
	statusLock sync.RWMutex
}

//初始化巡检器
func NewInspector() *Inspector {
	uuid, _ := uuid.NewV4()
	return &Inspector{
		ServerName:    uuid.String(),
		EtcdClient:    utils.GetEtcdClient(),
		ExecStatus:    true,
		InspectTicker: time.NewTicker(5 * time.Second),
		statusLock:    sync.RWMutex{},
	}
}

//启动巡检器
func (ict *Inspector) InspctorStart(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Inspector start ...")

	//启动监控最后一次轮巡检时间戳
	wg.Add(1)
	go ict.WatchInspectLastTime(wg)

	//启动定时监控是否可执行状态
	wg.Add(1)
	go ict.CheckExecStatus(wg)
	wg.Wait()

}
func (ict *Inspector) SetExecStatus(status bool) {
	ict.statusLock.Lock()
	defer ict.statusLock.Unlock()
	ict.ExecStatus = status

}

// watch etcd中的巡检标志，标志为最后一次巡检的时间戳，若收到其他服务器端巡检时间戳则更改 ExecStatus 为false
func (ict *Inspector) WatchInspectLastTime(wg *sync.WaitGroup) {
	defer wg.Done()
	rch := ict.EtcdClient.Watch(context.Background(), LastInspectionTime) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			ict.SetExecStatus(false)
		}
	}
}

//巡检轮询器，定时轮训巡检标志。
//若ExecStatus为true 则向etcd发送当前时间戳，并执行巡检过程
func (ict *Inspector) CheckExecStatus(wg *sync.WaitGroup) {
	defer wg.Done()
	for range ict.InspectTicker.C {

		if ict.ExecStatus == true {
			unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
			if _, err := ict.EtcdClient.Put(context.Background(), LastInspectionTime, unixNano); err != nil {
				fmt.Println(err)
				continue
			}
			ict.execInspect()
		} else {
			ict.SetExecStatus(true)
		}
	}
}

//巡检函数，检查etcd分布式锁，若有任务执行则向LastInspectionTime发送时间戳，若没有锁则上锁执行任务
func (ict *Inspector) execInspect() {
	//利用分布式锁避免执行冲突
	//有其他任务执行时，发送LastInspectionTime
	session, err := concurrency.NewSession(ict.EtcdClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := concurrency.NewMutex(session, InspectLockKey)
	unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := m.Lock(context.TODO()); err != nil {
		if _, err := ict.EtcdClient.Put(context.Background(), LastInspectionTime, unixNano); err != nil {
			fmt.Println(err)
			return
		}
	}

	//ToDo 巡检逻辑
	fmt.Println("Execute inspection task: ", ict.ServerName)
	time.Sleep(30 * time.Second)

	defer m.Unlock(context.TODO())
}
