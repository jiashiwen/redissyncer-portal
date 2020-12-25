// Package inspection 实现定期巡检
// 实现方式：
//     服务器 watch LastInspectionTime key
//     当收到 LastInspectionTime 变更则修改是否执行的标志为false
//     轮训器定时检查 execstatus flag，若值为true则触发巡检，若为false则修改值为true
//     巡检执行前开启分布式锁，任务执行完成解锁
//     当巡检器加锁发生互斥时，修改LastInspectionTime后退出

package inspection

import "sync"

type Inspect interface {
	Start(wg sync.WaitGroup)
	Stop()
}
