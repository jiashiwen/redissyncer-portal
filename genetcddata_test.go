package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"redissyncer-portal/commons"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"redissyncer-portal/node"
	"redissyncer-portal/utils"
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
				clientv3.OpPut(global.TasksTaskIDPrefix+v+strconv.Itoa(i), "{}"),
				//put  TasksNode
				clientv3.OpPut(global.TasksNodePrefix+strconv.Itoa(k)+"/"+v+strconv.Itoa(i), "{}"),
				//put TasksGroupId
				clientv3.OpPut(global.TasksGroupIDPrefix+v+strconv.Itoa(i), v+strconv.Itoa(i)),
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

	global.RSPViper = core.Viper()
	global.RSPLog = core.Zap()

	etcdClient := global.GetEtcdClient()
	defer etcdClient.Close()
	kv := clientv3.NewKV(etcdClient)

	nameSeed := commons.RandString(4)

	for i := 0; i < 50; i++ {
		taskstatus := GenTaskStatusData(nameSeed + strconv.Itoa(i))
		taskstatus.NodeID = strconv.Itoa(i)
		nodeskey := global.NodesPrefix + global.NodeTypeRedissyncer + "/" + strconv.Itoa(i)
		nodestatus := node.NodeStatus{
			//节点类型
			NodeType: global.NodeTypeRedissyncer,

			//节点Id
			NodeID: strconv.Itoa(i),

			//节点ip地址
			NodeAddr: "127.0.0.1",

			//探活port
			NodePort: 8888,

			//探活url
			HeartbeatUrl: "/health",

			//是否在线
			Online: true,

			//最后上报时间，unix时间戳
			LastReportTime: 1615791131672,
		}
		nodestatusval, _ := json.Marshal(nodestatus)

		tasksTaskidKey := global.TasksTaskIDPrefix + taskstatus.TaskID
		tasksTaskidVal, _ := json.Marshal(taskstatus)
		tasksGroupidkey := global.TasksGroupIDPrefix + taskstatus.GroupID + "/" + taskstatus.TaskID
		tasksGroupidMap := map[string]string{"groupId": taskstatus.GroupID, "taskId": taskstatus.TaskID}
		tasksGroupidJSON, _ := json.Marshal(tasksGroupidMap)
		tasksStatuskey := global.TasksStatusPrefix + strconv.Itoa(taskstatus.Status) + "/" + taskstatus.TaskID
		tasksStatusVal := map[string]string{"taskId": taskstatus.TaskID}
		tasksStatusJSON, _ := json.Marshal(tasksStatusVal)
		tasksNameKey := global.TasksNamePrefix + taskstatus.TaskName
		// tasksOffsetKey := global.TasksOffsetPrefix + taskstatus.TaskID
		// tasksOffsetVal := global.TasksOffset{ReplID: taskstatus.ReplID, ReplOffset: 11223344}
		// tasksOffsetJSON, _ := json.Marshal(tasksOffsetVal)
		tasksNodeKey := global.TasksNodePrefix + strconv.Itoa(i) + "/" + taskstatus.TaskID
		tasksNodeVal := map[string]string{"nodeId": strconv.Itoa(i), "taskId": taskstatus.TaskID}
		tasksNodeJSON, _ := json.Marshal(tasksNodeVal)
		tasksMD5Key := global.TasksMd5Prefix + taskstatus.MD5
		tasksMD5Val := map[string]string{
			"groupId": taskstatus.GroupID, "nodeId": strconv.Itoa(i), "taskId": taskstatus.TaskID,
		}
		tasksMD5JSON, _ := json.Marshal(tasksMD5Val)
		kv.Txn(context.TODO()).Then(
			//put TasksTaskId
			clientv3.OpPut(tasksTaskidKey, string(tasksTaskidVal)),
			//put  TasksNode
			clientv3.OpPut(tasksNodeKey, string(tasksNodeJSON)),
			//put TasksGroupId
			clientv3.OpPut(tasksGroupidkey, string(tasksGroupidJSON)),
			//put TasksName
			clientv3.OpPut(tasksNameKey, string(tasksStatusJSON)),
			//put TasksStatus
			clientv3.OpPut(tasksStatuskey, string(tasksStatusJSON)),
			//put TasksMD5
			clientv3.OpPut(tasksMD5Key, string(tasksMD5JSON)),
			//put NodesRedissyncer
			clientv3.OpPut(nodeskey, string(nodestatusval)),
		).Commit()

		fmt.Println("key: ", tasksTaskidKey)

	}

}

func GenTaskStatusData(taskname string) *global.TaskStatus {
	taskid := uuid.Must(uuid.NewV4(), nil)

	sourceip := utils.RandomIp()
	targetip := utils.RandomIp()
	targeturi := []string{}
	targeturi = append(targeturi, "redis://"+targetip+":6379")
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
		SourceHost:         sourceip,
		SourcePassword:     "redistest0102",
		SourcePort:         6379,
		SourceRedisAddress: sourceip + ":6379",
		SourceRedisType:    1,
		SourceURI:          "redis://" + sourceip + ":6379?authPassword=redistest0102",
		SourceUserName:     "",
		Status:             int(global.TaskStatusTypeRDBRUNING),
		SyncType:           1,
		TargetACL:          false,
		TargetHost:         targetip,
		TargetPassword:     "",
		TargetPort:         6379,
		TargetRedisAddress: targetip + ":6379",
		TargetRedisType:    1,
		TargetURI:          targeturi,
		TargetUserName:     "",
		TaskID:             taskid.String(),
		TaskMsg:            "全量同步开始[同步任务启动]",
		TaskName:           taskname,
		TaskType:           1,
		TimeDeviation:      0,
	}
	taskStatus.MD5 = GetTaskMd5(&taskStatus)
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
