package service

import (
	"context"
	"encoding/json"
	"errors"
	"redissyncer-portal/global"
	"redissyncer-portal/httpquerry"
	"redissyncer-portal/node"
	"strconv"
	"strings"

	"github.com/coreos/etcd/clientv3"
)

func CreateTask(body string) (string, error) {
	// 节点选择 pairelist 给出节点map[nodeid]任务数量，按任务数量排序
	selector := node.NodeSelector{
		EtcdClient: global.GetEtcdClient(),
	}

	pairelist, err := selector.SelectNode()

	if err != nil {
		return err.Error(), err
	}

	// 按顺序检查节点可用后发送创建任务请求，若第一个节点不可以，顺序执行第二节点
	for _, v := range *pairelist {
		//获取节点ip、port
		getresp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix+global.NodeTypeRedissyncer+"/"+v.Key)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			continue
		}

		var node node.Node
		json.Unmarshal(getresp.Kvs[0].Value, &node)

		if httpquerry.NodeAlive(node.NodeAddr, strconv.Itoa(node.NodePort)) {
			//向选中的redissyncer-server发送创建任务请求
			req := httpquerry.New("http://" + node.NodeAddr + ":" + strconv.Itoa(node.NodePort))
			req.Api = httpquerry.UrlCreateTask
			req.Body = body
			global.RSPLog.Sugar().Debug("exec create task")
			return req.ExecRequest()
		}
		continue
	}
	return "", errors.New("no node selected")
}

//Start task
func StartTask(body string) (string, error) {

	return "", nil

}

//Stop task by task ids
func StopTaskByIds(ids []string) (string, error) {
	return "", nil
}

//RemoveTasks
func RemoveTask(taskID string) error {

	var taskStatus global.TaskStatus
	//获取节点status
	statusResp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+taskID)
	if err != nil {
		global.RSPLog.Sugar().Debug(err)
		return err
	}

	//判断taskid是否存在
	if len(statusResp.Kvs) == 0 {
		return errors.New("taskid not exists")
	}

	//判断任务是否为停止状态
	json.Unmarshal(statusResp.Kvs[0].Value, &taskStatus)
	if taskStatus.Status != int(global.TaskStatusTypeSTOP) &&
		taskStatus.Status != int(global.TaskStatusTypeBROKEN) &&
		taskStatus.Status != int(global.TaskStatusTypeFINISH) {
		return errors.New("task is running")

	}

	//清理任务数据
	kv := clientv3.NewKV(global.GetEtcdClient())
	txnResp, err := kv.Txn(context.TODO()).Then(
		//del TasksTaskId
		clientv3.OpDelete(global.TasksTaskIDPrefix+taskStatus.TaskID),
		//del  TasksNode
		clientv3.OpDelete(global.TasksNodePrefix+taskStatus.NodeID+"/"+taskStatus.TaskID),
		//del TasksGroupId
		clientv3.OpDelete(global.TasksGroupIDPrefix+taskStatus.GroupID+"/"+taskStatus.TaskID),
		//del TasksStatus
		clientv3.OpDelete(global.TasksStatusPrefix+strconv.Itoa(taskStatus.Status)+"/"+taskStatus.TaskID),
		//del TasksOffset
		clientv3.OpDelete(global.TasksOffsetPrefix+taskStatus.TaskID),
		//del TasksName
		clientv3.OpDelete(global.TasksNamePrefix+taskStatus.TaskName),
		//del TasksType
		clientv3.OpDelete(global.TasksTypePrefix+strconv.Itoa(taskStatus.TaskType)+"/"+taskStatus.TaskID),
		//del TasksBigkey
		clientv3.OpDelete(global.TasksBigkeyPrefix+taskStatus.TaskID, clientv3.WithPrefix()),
		//del TasksMd5
		clientv3.OpDelete(global.TasksMd5Prefix+taskStatus.MD5),
	).Commit()

	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return errors.New("txn commit fail")
	}

	return nil

}

//Remove task by name
func RemoveTaskByName(taskName string) (string, error) {
	return "", nil

}

// GetTaskStatus 获取同步任务状态
func GetTaskStatus(ids []string) ([]*global.TaskStatus, error) {
	var tasksStatus []*global.TaskStatus
	for _, id := range ids {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+id)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskStatus := global.TaskStatus{}
				json.Unmarshal(v.Value, &taskStatus)
				tasksStatus = append(tasksStatus, &taskStatus)
			}
		}
	}

	return tasksStatus, nil
}

// GetTaskStatusByName 根据名字查找任务状态
func GetTaskStatusByName(taskNames []string) ([]*global.TaskStatus, error) {

	var taskIds []string
	for _, name := range taskNames {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksNamePrefix+name)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskStatus := global.TaskStatus{}
				json.Unmarshal(v.Value, &taskStatus)
				taskIds = append(taskIds, taskStatus.TaskID)
			}
		}
	}

	return GetTaskStatus(taskIds)
}

// GetTaskStatusByGroupID 根据groupid获取任务状态
func GetTaskStatusByGroupID(groupIDs []string) ([]*global.TaskStatus, error) {
	taskIDsArry := []string{}
	for _, groupID := range groupIDs {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksGroupIDPrefix+groupID, clientv3.WithPrefix())
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskIDsArry = append(taskIDsArry, strings.Split(string(v.Key), "/")[4])
			}
		}

	}

	return GetTaskStatus(taskIDsArry)
}

// @title    GetSameTaskNameIds
// @description   获取同名任务列表
// @auth      Jsw             2020/7/1   10:57
// @param    taskName        string         "任务名称"
// @return    taskIds        []string         "任务id数组"
func GetSameTaskNameIds(taskName string) ([]string, error) {

	// var existTaskIds []string
	// listJsonMap := make(map[string]interface{})
	// listJsonMap["regulation"] = "bynames"
	// listJsonMap["tasknames"] = strings.Split(taskName, ",")
	// listJsonStr, err := json.Marshal(listJsonMap)
	// if err != nil {
	// 	return nil, err
	// }

	// r.Api = UrlListTasks
	// r.Body = string(listJsonStr)

	// listResp, err := r.ExecRequest()
	// if err != nil {
	// 	return nil, err
	// }
	// taskList := gjson.Get(listResp, "data").Array()

	// if len(taskList) > 0 {
	// 	for _, v := range taskList {
	// 		existTaskIds = append(existTaskIds, gjson.Get(v.String(), "taskId").String())
	// 	}
	// }
	// return existTaskIds, nil
	return nil, nil
}
