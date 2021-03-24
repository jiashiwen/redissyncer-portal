package service

import (
	"context"
	"encoding/json"
	"errors"
	"redissyncer-portal/global"
	"redissyncer-portal/httpquerry"
	"redissyncer-portal/httpserver/model"
	"redissyncer-portal/httpserver/model/response"
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
		getResp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix+global.NodeTypeRedissyncer+"/"+v.Key)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			continue
		}

		var node node.Node
		json.Unmarshal(getResp.Kvs[0].Value, &node)

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
func StartTask(body model.TaskStartBody) (string, error) {

	var taskStatus global.TaskStatus
	//通过taskID 获取TaskStatus
	statusResp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+body.TaskID)
	if err != nil {
		global.RSPLog.Sugar().Debug(err)
		return "", err
	}

	//判断taskid是否存在
	if len(statusResp.Kvs) == 0 {
		return "", errors.New("taskid not exists")
	}

	//判断任务是否为停止状态
	json.Unmarshal(statusResp.Kvs[0].Value, &taskStatus)
	if taskStatus.Status != int(global.TaskStatusTypeSTOP) &&
		taskStatus.Status != int(global.TaskStatusTypeBROKEN) &&
		taskStatus.Status != int(global.TaskStatusTypeFINISH) {
		return "", errors.New("task is running")

	}

	//获取节点ip、port
	nodeResp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix+global.NodeTypeRedissyncer+"/"+taskStatus.NodeID)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		//return "", errors.New("node not exist")
		return "", err
	}

	if len(nodeResp.Kvs) == 0 {
		return "", errors.New("node not exist")
	}

	var node node.NodeStatus
	json.Unmarshal(nodeResp.Kvs[0].Value, &node)

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return "", err
	}

	//拼装httpqurerry
	req := httpquerry.New("http://" + node.NodeAddr + ":" + strconv.Itoa(node.NodePort))
	req.Api = httpquerry.UrlStartTask
	req.Body = string(bodyJSON)

	return req.ExecRequest()

}

//Stop task by task ids
func StopTaskById(taskID string) (string, error) {
	var taskStatus global.TaskStatus
	//通过taskID 获取TaskStatus
	statusResp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+taskID)
	if err != nil {
		global.RSPLog.Sugar().Debug(err)
		return "", err
	}

	//判断taskid是否存在
	if len(statusResp.Kvs) == 0 {
		return "", errors.New("taskid not exists")
	}

	//判断任务是否为停止状态
	json.Unmarshal(statusResp.Kvs[0].Value, &taskStatus)
	if taskStatus.Status == int(global.TaskStatusTypeSTOP) ||
		taskStatus.Status == int(global.TaskStatusTypeBROKEN) ||
		taskStatus.Status == int(global.TaskStatusTypeFINISH) {
		return "", errors.New("task is stopped")

	}

	//获取节点ip、port
	nodeResp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix+global.NodeTypeRedissyncer+"/"+taskStatus.NodeID)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return "", err
	}

	if len(nodeResp.Kvs) == 0 {
		return "", errors.New("node not exist")
	}

	var node node.NodeStatus
	taskStopBody := model.TaskStopBodyToNode{
		TaskIDs: []string{taskStatus.TaskID},
	}
	json.Unmarshal(nodeResp.Kvs[0].Value, &node)

	bodyJSON, err := json.Marshal(taskStopBody)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return "", err
	}

	//拼装httpqurerry
	req := httpquerry.New("http://" + node.NodeAddr + ":" + strconv.Itoa(node.NodePort))
	req.Api = httpquerry.UrlStopTask
	req.Body = string(bodyJSON)

	return req.ExecRequest()

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

//获取任务状态
func GetTaskStatus(id string) (*global.TaskStatus, error) {

	resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+id)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		err := errors.New("task not exist")
		global.RSPLog.Sugar().Error(err)
		return nil, err
	}

	taskStatus := global.TaskStatus{}
	json.Unmarshal(resp.Kvs[0].Value, &taskStatus)

	return &taskStatus, nil
}

// GetTaskStatusByIDs 获取同步任务状态
func GetTaskStatusByIDs(ids []string) []*response.TaskStatusResult {
	var taskStatusResultArray []*response.TaskStatusResult
	for _, id := range ids {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix+id)
		if err != nil {
			errorResult := response.ErrorResult{
				Code: global.ErrorSystemError,
				Msg:  err.Error(),
			}
			taskStatusResult := response.TaskStatusResult{
				TaskID:     id,
				Errors:     &errorResult,
				TaskStatus: nil,
			}
			taskStatusResultArray = append(taskStatusResultArray, &taskStatusResult)
			global.RSPLog.Sugar().Error(err)
			continue
		}

		if len(resp.Kvs) == 0 {
			errorResult := response.ErrorResult{
				Code: global.ErrorTaskNotExists,
				Msg:  global.ErrorTaskNotExists.String(),
			}
			taskStatusResult := response.TaskStatusResult{
				TaskID:     id,
				Errors:     &errorResult,
				TaskStatus: nil,
			}
			taskStatusResultArray = append(taskStatusResultArray, &taskStatusResult)
			continue
		}

		taskStatus := global.TaskStatus{}

		if err := json.Unmarshal(resp.Kvs[0].Value, &taskStatus); err != nil {
			errorResult := response.ErrorResult{
				Code: global.ErrorSystemError,
				Msg:  err.Error(),
			}
			taskStatusResult := response.TaskStatusResult{
				TaskID:     id,
				Errors:     &errorResult,
				TaskStatus: nil,
			}
			taskStatusResultArray = append(taskStatusResultArray, &taskStatusResult)
			global.RSPLog.Sugar().Error(err)
			continue
		}

		taskStatusResult := response.TaskStatusResult{
			TaskID:     id,
			Errors:     nil,
			TaskStatus: &taskStatus,
		}
		taskStatusResultArray = append(taskStatusResultArray, &taskStatusResult)

	}

	return taskStatusResultArray
}

// GetTaskStatusByName 根据名字查找任务状态
func GetTaskStatusByName(taskNames []string) []*response.TaskStatusResultByName {
	var taskStatusByNameArray []*response.TaskStatusResultByName
	var taskIds []string
	for _, name := range taskNames {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksNamePrefix+name)
		if err != nil {
			errorCode := response.ErrorResult{
				Code: global.ErrorSystemError,
				Msg:  err.Error(),
			}
			taskStatusByName := response.TaskStatusResultByName{
				TaskName:   name,
				Errors:     &errorCode,
				TaskStatus: nil,
			}
			taskStatusByNameArray = append(taskStatusByNameArray, &taskStatusByName)
			global.RSPLog.Sugar().Error(err)
			continue
		}

		if len(resp.Kvs) == 0 {
			errorCode := response.ErrorResult{
				Code: global.ErrorTaskNotExists,
				Msg:  global.ErrorTaskNotExists.String(),
			}
			taskStatusByName := response.TaskStatusResultByName{
				TaskName:   name,
				Errors:     &errorCode,
				TaskStatus: nil,
			}
			taskStatusByNameArray = append(taskStatusByNameArray, &taskStatusByName)
			continue
		}

		for _, v := range resp.Kvs {
			taskStatus := global.TaskStatus{}
			json.Unmarshal(v.Value, &taskStatus)
			taskIds = append(taskIds, taskStatus.TaskID)
		}
	}

	for _, v := range GetTaskStatusByIDs(taskIds) {
		taskStatusByName := response.TaskStatusResultByName{
			TaskName:   v.TaskStatus.TaskName,
			Errors:     v.Errors,
			TaskStatus: v.TaskStatus,
		}
		taskStatusByNameArray = append(taskStatusByNameArray, &taskStatusByName)
	}
	return taskStatusByNameArray
}

// GetTaskStatusByGroupIDs 根据groupid获取任务状态
func GetTaskStatusByGroupIDs(groupIDs []string) []*response.TaskStatusResultByGroupID {
	var taskStatusResultByGroupIDArray []*response.TaskStatusResultByGroupID
	var groupIDTaskIDMap map[string][]string
	for _, groupID := range groupIDs {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksGroupIDPrefix+groupID, clientv3.WithPrefix())
		if err != nil {
			errorCode := response.ErrorResult{
				Code: global.ErrorSystemError,
				Msg:  err.Error(),
			}
			taskStatusResultByGroupID := response.TaskStatusResultByGroupID{
				GroupID:         groupID,
				Errors:          &errorCode,
				TaskStatusArray: nil,
			}
			taskStatusResultByGroupIDArray = append(taskStatusResultByGroupIDArray, &taskStatusResultByGroupID)
			global.RSPLog.Sugar().Error(err)
			continue
		}

		if len(resp.Kvs) == 0 {
			errorCode := response.ErrorResult{
				Code: global.ErrorTaskGroupNotExists,
				Msg:  global.ErrorTaskGroupNotExists.String(),
			}
			taskStatusResultByGroupID := response.TaskStatusResultByGroupID{
				GroupID:         groupID,
				Errors:          &errorCode,
				TaskStatusArray: nil,
			}
			taskStatusResultByGroupIDArray = append(taskStatusResultByGroupIDArray, &taskStatusResultByGroupID)
			global.RSPLog.Sugar().Error(err)
			continue
		}
		var taskIDsArray []string
		for _, v := range resp.Kvs {
			taskIDsArray = append(taskIDsArray, strings.Split(string(v.Key), "/")[4])
		}

		groupIDTaskIDMap[groupID] = taskIDsArray
	}

	for k, v := range groupIDTaskIDMap {
		taskStatusByGroupID := response.TaskStatusResultByGroupID{
			GroupID:         k,
			Errors:          nil,
			TaskStatusArray: GetTaskStatusByIDs(v),
		}
		taskStatusResultByGroupIDArray = append(taskStatusResultByGroupIDArray, &taskStatusByGroupID)

	}
	return taskStatusResultByGroupIDArray
}

func GetAllTaskStatus(queryid string) []*global.TaskStatus {
	
	return nil
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