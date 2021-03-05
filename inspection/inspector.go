package inspection

import (
	"context"
	"encoding/json"
	"redissyncer-portal/global"
	"redissyncer-portal/httpquerry"
	"redissyncer-portal/node"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gofrs/uuid"
)

const (
	// LastInspectionTime 最后一次的巡检时间戳
	LastInspectionTime = "/inspect/lastinspectiontime"

	// InspectLockKey 巡检锁
	InspectLockKey = "/inspect/execlock"
)

//Inspector 巡检器
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

	//巡检器context
	InspectorContext context.Context

	//巡检器cancle
	InspectorCancel context.CancelFunc
}

// NewInspector 初始化巡检器
func NewInspector() *Inspector {
	name, _ := uuid.NewV4()
	ctx, cancel := context.WithCancel(context.Background())
	return &Inspector{
		ServerName:       name.String(),
		EtcdClient:       global.GetEtcdClient(),
		ExecStatus:       true,
		InspectTicker:    time.NewTicker(5 * time.Second),
		statusLock:       sync.RWMutex{},
		InspectorContext: ctx,
		InspectorCancel:  cancel,
	}
}

//Start 启动巡检器
func (ict *Inspector) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	defer ict.InspectorCancel()
	global.RSPLog.Sugar().Info("Inspector start ...")

	//启动监控最后一次轮巡检时间戳
	wg.Add(1)
	go ict.WatchInspectLastTime(ict.InspectorContext, wg)

	//启动定时监控是否可执行状态
	wg.Add(1)
	go ict.CheckExecStatus(ict.InspectorContext, wg)
	wg.Wait()

}

// Stop 停止检查器
func (ict *Inspector) Stop() {
	ict.InspectorCancel()
}

// SetExecStatus 设置执行状态来判断当 ticker 触发时是否执行巡检
func (ict *Inspector) SetExecStatus(status bool) {
	ict.statusLock.Lock()
	defer ict.statusLock.Unlock()
	ict.ExecStatus = status

}

// WatchInspectLastTime watch etcd中的巡检标志，标志为最后一次巡检的时间戳，若收到其他服务器端巡检时间戳则更改 ExecStatus 为false
func (ict *Inspector) WatchInspectLastTime(ctx context.Context, wg *sync.WaitGroup) {
	//func (ict *Inspector) WatchInspectLastTime(ctx context.Context) {
	defer wg.Done()

	rch := ict.EtcdClient.Watch(context.Background(), LastInspectionTime) // <-chan WatchResponse

	for {
		select {
		case resp, _ := <-rch:
			for _, ev := range resp.Events {
				if string(ev.Kv.Key) == LastInspectionTime {
					global.RSPLog.Sugar().Debugf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					ict.SetExecStatus(false)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// CheckExecStatus 巡检轮询器，定时轮训巡检标志。
//若ExecStatus为true 则向etcd发送当前时间戳，并执行巡检过程
func (ict *Inspector) CheckExecStatus(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for range ict.InspectTicker.C {
		select {
		case <-ctx.Done():
			return
		default:
			if ict.ExecStatus == true {
				unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
				if _, err := ict.EtcdClient.Put(context.Background(), LastInspectionTime, unixNano); err != nil {
					global.RSPLog.Sugar().Error(err)
					continue
				}
				ict.execInspect()
			} else {
				ict.SetExecStatus(true)
			}
		}
	}
}

//巡检函数，检查etcd分布式锁，若有任务执行则向LastInspectionTime发送时间戳，若没有锁则上锁执行任务
func (ict *Inspector) execInspect() {
	//利用分布式锁避免执行冲突
	//有其他任务执行时，发送LastInspectionTime
	session, err := concurrency.NewSession(ict.EtcdClient)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return
	}

	m := concurrency.NewMutex(session, InspectLockKey)
	unixNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := m.Lock(context.TODO()); err != nil {
		if _, err := ict.EtcdClient.Put(context.Background(), LastInspectionTime, unixNano); err != nil {
			global.RSPLog.Sugar().Error(err)
			return
		}
	}

	// 巡检逻辑
	if err := ict.nodeHealthCheck(); err != nil {
		global.RSPLog.Sugar().Error(err)
	}

	if err := m.Unlock(context.TODO()); err != nil {
		global.RSPLog.Sugar().Error(err)
	}
}

//节点健康检查
//取LastReportTime字段，与当前时间戳对比若差值大于阈值则说明节点可能离线，检查节点health情况若宕机则改写Online为false
func (ict *Inspector) nodeHealthCheck() error {
	getResp, err := ict.EtcdClient.Get(context.Background(), "/nodes", clientv3.WithPrefix())

	if err != nil {
		return err
	}
	var nodeStatus node.NodeStatus
	for _, v := range getResp.Kvs {
		if err := json.Unmarshal([]byte(v.Value), &nodeStatus); err != nil {
			return err
		}
		//ToDo 健康检查逻辑
		//本地时间戳（毫秒）
		localUnixTimestamp := time.Now().UnixNano() / 1e6
		if localUnixTimestamp-nodeStatus.LastReportTime > 10000 {
			if nodeStatus.Online == false {
				return nil
			}
			//执行探活,若确定node离线则修改node online属性为false
			if !httpquerry.NodeAlive(nodeStatus.NodeAddr, strconv.Itoa(nodeStatus.NodePort)) {
				nodeStatus.Online = false
				statusJSON, err := json.Marshal(&nodeStatus)
				if err != nil {
					global.RSPLog.Sugar().Error(err)
					return err
				}
				if _, err := global.GetEtcdClient().Put(context.Background(), string(v.Key), string(statusJSON)); err != nil {
					global.RSPLog.Sugar().Error(err)
					return err
				}
			}
		}
	}
	return nil
}
