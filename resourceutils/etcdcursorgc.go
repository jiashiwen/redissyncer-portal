package resourceutils

import (
	"context"
	"errors"
	"redissyncer-portal/commons"
	"redissyncer-portal/global"
	"sync"
	"time"
)

var cursorGC EtcdCursorGC

type EtcdCursorGC struct {
	GCTicker *time.Ticker //GC定时触发
	TimeDiff int64        //当前时间与最后访问时间的差值，当大于差值时执行删除操作
	Ctx      context.Context
	GCCancel context.CancelFunc //GC Cancle function
}

func NewEtcdCursorGC(tickerTime int, timeDiff int64) {
	cursorGC = EtcdCursorGC{
		TimeDiff: 60000,
		GCTicker: time.NewTicker(time.Duration(tickerTime) * time.Second),
	}

	if timeDiff > 1000 {
		cursorGC.TimeDiff = timeDiff
	}
}

func (gc *EtcdCursorGC) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(context.Background())
	gc.GCCancel = cancel
	gc.Ctx = ctx
	wg.Add(1)
	go gc.DoGC(ctx, wg)
}

func (gc *EtcdCursorGC) Stop() {
	gc.GCCancel()
}

//清理本地cursorMap
func (gc *EtcdCursorGC) DoGC(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for range gc.GCTicker.C {
		select {
		case <-ctx.Done():
			return
		default:
			cursorQueryMap := *GetCursorQueryMap()
			currentUnixTime := commons.GetCurrentUnixMillisecond()
			for _, v := range cursorQueryMap {
				if currentUnixTime-v.LastQueryTimeStamp > gc.TimeDiff {
					delete(cursorQueryMap, v.QueryID)
				}
			}
			global.RSPLog.Sugar().Debug("current unix millisecond is: ", currentUnixTime)
		}
	}
}

//初始化cursorGC
func InitCursorGC(gc EtcdCursorGC) {
	cursorGC = gc
}

//启动GC
func StartCursorGC(wg *sync.WaitGroup) {
	defer wg.Done()
	if commons.IsNil(cursorGC) {
		panic(errors.New("cursor GC not init"))
	}
	wg.Add(1)
	go cursorGC.Start(wg)
	global.RSPLog.Sugar().Info("cursor GC start ...")
}

//停止GC
func StopCursorGC() {
	cursorGC.Stop()
}

//变更GC
