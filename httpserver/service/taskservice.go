package service

import (
	"context"
	"redissyncer-portal/global"
	"redissyncer-portal/node"
)

func CreateTask(body string) error {
	// 节点选择 pairelist 给出节点map[nodeid]任务数量，按任务数量排序
	selector := node.NodeSelector{
		EtcdClient: global.GetEtcdClient(),
	}

	pairelist, err := selector.SelectNode()

	if err != nil {
		return err
	}

	// 按顺序检查节点可用后发送创建任务请求，若第一个节点不可以，顺序执行第二节点
	for _, v := range *pairelist {
		global.RSPLog.Sugar().Debug("nodeId:", v.Key, "tasks:", v.Value)
	}
	return nil
}

//Start task
func StartTask(taskid string) (string, error) {

	return "", nil

}

//Stop task by task ids
func StopTaskByIds(ids []string) (string, error) {
	return "", nil
}

//Remove task by name
func RemoveTaskByName(taskName string) (string, error) {
	return "", nil

}

// GetTaskStatus 获取同步任务状态
func GetTaskStatus(ids []string) (*map[string]string, error) {
	taskStatusMap := make(map[string]string)
	for _, id := range ids {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskidPrefix+id)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskStatusMap[string(v.Key)] = string(v.Value)
			}
		}
	}

	return &taskStatusMap, nil
}

// GetTaskStatusByName 根据名字查找任务状态
func GetTaskStatusByName(taskNames []string) (*map[string]string, error) {
	taskStatusMap := make(map[string]string)
	for _, name := range taskNames {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksNamePrefix+name)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskStatusMap[string(v.Key)] = string(v.Value)
			}
		}
	}

	return &taskStatusMap, nil
}

// GetTaskStatusByGroupID 根据groupid获取任务状态
func GetTaskStatusByGroupID(groupIDs []string) (*map[string]string, error) {
	taskStatusMap := make(map[string]string)
	for _, groupID := range groupIDs {
		resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksGroupidPrefix+groupID)
		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return nil, err
		}
		if len(resp.Kvs) > 0 {
			for _, v := range resp.Kvs {
				taskStatusMap[string(v.Key)] = string(v.Value)
			}
		}
	}

	return &taskStatusMap, nil
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
