package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	uuid "github.com/satori/go.uuid"
)

// const (
// 	IDSeed           = "/uniqid/idseed"
// 	TasksTaskID      = "/tasks/taskid"
// 	TasksNode        = "/tasks/node/"
// 	TasksGroupID     = "/tasks/groupid/"
// 	TasksStatus      = "/tasks/status/"
// 	TasksName        = "/tasks/name/"
// 	NodesRedissyncer = "/nodes/redissyncerserver/"
// )

var lock sync.Mutex

func TestGenEtcdData(t *testing.T) {
	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

	etcdClient := global.GetEtcdClient()
	defer etcdClient.Close()

	//获取uniqid，加锁
	session, err := concurrency.NewSession(etcdClient)
	if err != nil {
		// global.RSPLog.Sugar().Error(err)
		t.Error(err)
		return
	}

	m := concurrency.NewMutex(session, global.IDSeed)
	// unixTimeStamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

	if err := m.Lock(context.TODO()); err != nil {
		t.Error(err)
		return
	}

	getResp, err := etcdClient.Get(context.Background(), global.IDSeed)
	if err != nil {
		t.Error(err)
		return
	}

	//判断idseed是否存在，若存在值+1，若不存在value=1
	if len(getResp.Kvs) == 0 {
		etcdClient.Put(context.Background(), global.IDSeed, "1")
	} else {
		value := getResp.Kvs[0].Value
		valueInt64, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			t.Error(err)
			return
		}
		etcdClient.Put(context.Background(), global.IDSeed, strconv.FormatInt(valueInt64+int64(1), 10))
	}

	taskids := GenGlobalUniqID(string(getResp.Kvs[0].Value), 10)
	if err := m.Unlock(context.TODO()); err != nil {
		global.RSPLog.Sugar().Error(err)
	}

	kv := clientv3.NewKV(etcdClient)

	//生成模拟tasks
	for k, v := range taskids {
		fmt.Println(v)

		for i := 0; i < k; i++ {
			kv.Txn(context.TODO()).Then(

				//put NodesRedissyncer
				clientv3.OpPut(global.NodeTypeRedissyncer+strconv.Itoa(i), "{}"),
				//put TasksTaskId
				clientv3.OpPut(global.TasksTaskidPrefix+v+strconv.Itoa(i), "{}"),
				//put  TasksNode
				clientv3.OpPut(global.TasksNodePrefix+strconv.Itoa(k)+"/"+v+strconv.Itoa(i), "{}"),
				//put TasksGroupId
				clientv3.OpPut(global.TasksGroupidPrefix+v+strconv.Itoa(i), v+strconv.Itoa(i)),
				//put TasksName
				clientv3.OpPut(global.TasksNamePrefix+v+strconv.Itoa(i), v+strconv.Itoa(i)),
				//put TasksStatus
				clientv3.OpPut(global.TasksStatusPrefix+strconv.Itoa(int(global.TaskStatusTypeCREATING))+"/"+v+strconv.Itoa(i), v+strconv.Itoa(i)),
			).Commit()
		}

	}
}

func TestParseStatusToStruct(t *testing.T) {
	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

	etcdClient := global.GetEtcdClient()
	defer etcdClient.Close()

	resp, _ := etcdClient.Get(context.Background(), "/tasks/taskid/DA38A1FB121942FC90432BDA56572352")
	var status global.TaskStatus

	json.Unmarshal(resp.Kvs[0].Value, &status)

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, resp.Kvs[0].Value, "", "\t")

	t.Log(string(resp.Kvs[0].Value))
	t.Log(string(prettyJSON.Bytes()))

	t.Logf("%+v", status)
	t.Logf("%d", global.TaskStatusTypeBROKEN)

}

func TestGenTaskDataForEtcd(t *testing.T) {

	GenTaskStatusData("abc")
}

func GenTaskStatusData(taskname string) *global.TaskStatus {
	taskid := uuid.Must(uuid.NewV4(), nil)

	taskStatus := global.TaskStatus{
		Afresh:             true,
		AllKeyCount:        286,
		AutoStart:          false,
		BatchSize:          500,
		CommandFilter:      "SET,DEL,FLUSHALL",
		DBMapper:           "{}",
		ErrorCount:         1,
		ExpandJSON:         "{\"brokenReason\":\"\"}",
		FileAddress:        "",
		FilterType:         "NONE",
		GroupID:            taskid.String(),
		ID:                 taskid.String(),
		KeyFilter:          "",
		LastKeyCommitTime:  1615451229743,
		LastKeyUpdateTime:  1615451232404,
		MD5:                "cb654f9bd7e554416eff90c8f0f6a047",
		Offset:             -1,
		OffsetPlace:        1,
		RdbKeyCount:        7489,
		RdbVersion:         7,
		RealKeyCount:       494,
		RedisVersion:       3.2,
		ReplID:             "5dea6a394f214eb29e3da2d282c5c744a06d8fe8",
		SourceACL:          false,
		SourceHost:         "116.196.82.161",
		SourcePassword:     "redistest0102",
		SourcePort:         6379,
		SourceRedisAddress: "116.196.82.161:6379",
		SourceRedisType:    1,
		SourceURI:          "redis://116.196.82.161:6379?authPassword=redistest0102",
		SourceUserName:     "",
		Status:             6,
		SyncType:           1,
		TargetACL:          false,
		TargetHost:         "127.0.0.1",
		TargetPassword:     "",
		TargetPort:         6379,
		TargetRedisAddress: "127.0.0.1:6379",
		TargetRedisType:    1,
		TargetURI:          []string{"redis://127.0.0.1:6379"},
		TargetUserName:     "",
		TaskID:             taskid.String(),
		TaskMsg:            "全量同步开始[同步任务启动]",
		TaskName:           taskname,
		TaskType:           1,
		TimeDeviation:      0,
	}
	taskStatus.MD5 = GetTaskMd5(&taskStatus)
	fmt.Println()
	fmt.Println(taskStatus.Status == 0)
	fmt.Printf("%+v", taskStatus)

	return &taskStatus

}

func GetTaskMd5(taskStatus *global.TaskStatus) string {

	taskMD5str := ""

	if taskStatus.TargetRedisAddress != "" {
		taskMD5str = taskMD5str + taskStatus.TargetRedisAddress + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	if taskStatus.TargetPassword != "" {
		taskMD5str = taskMD5str + taskStatus.TargetPassword + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	if taskStatus.SourceRedisAddress != "" {
		taskMD5str = taskMD5str + taskStatus.SourceRedisAddress + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	if taskStatus.SourcePassword != "" {
		taskMD5str = taskMD5str + taskStatus.SourcePassword + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	if taskStatus.FileAddress != "" {
		taskMD5str = taskMD5str + taskStatus.FileAddress + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	if taskStatus.TaskName != "" {
		taskMD5str = taskMD5str + taskStatus.TaskName + "_"
	} else {
		taskMD5str = taskMD5str + "null" + "_"
	}

	data := []byte(taskMD5str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func GenGlobalUniqID(idseed string, size int) []string {
	lock.Lock()
	defer lock.Unlock()
	unixTimeStamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	var ids []string

	for i := 0; i < size; i++ {
		ids = append(ids, idseed+"_"+unixTimeStamp+"_"+strconv.Itoa(i))
	}

	return ids

}
