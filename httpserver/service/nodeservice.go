package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	set "github.com/deckarep/golang-set"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver/model"
	"redissyncer-portal/httpserver/model/response"
	"redissyncer-portal/node"
	"redissyncer-portal/resourceutils"
	"strings"
)

//获取节点状态
func NodesStatus() ([]*response.NodeResult, error) {
	var resultArray []*response.NodeResult
	resp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, v := range resp.Kvs {
		var nodeStatus node.NodeStatus
		var nodeResult response.NodeResult
		if err := json.Unmarshal(v.Value, &nodeStatus); err != nil {
			nodeResult.Errors = &global.Error{
				Code: global.ErrorSystemError,
				Msg:  err.Error(),
			}
			resultArray = append(resultArray, &nodeResult)
			continue
		}
		nodeResult.NodeType = nodeStatus.NodeType
		nodeResult.NodeID = nodeStatus.NodeID
		nodeResult.NodeStatus = &nodeStatus
		resultArray = append(resultArray, &nodeResult)
	}
	return resultArray, nil
}

//获取所有节点类型
func NodeAllTypes() ([]string, error) {
	var types []string
	typesSet := set.NewSet()
	allNodesResp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, v := range allNodesResp.Kvs {
		typesSet.Add(strings.Split(string(v.Key), "/")[2])
	}

	fmt.Println(typesSet.ToSlice())
	for _, v := range typesSet.ToSlice() {
		types = append(types, fmt.Sprint(v))
	}
	return types, err
}

//删除节点
func RemoveNode(model model.RemoveNodeModel) *global.Error {
	var nodeStatus node.NodeStatus
	//判断节点是否为停止状态
	resp, err := global.GetEtcdClient().Get(context.Background(), global.NodesPrefix+model.NodeType+"/"+model.NodeID)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  err.Error(),
		}
	}

	if len(resp.Kvs) == 0 {
		return &global.Error{
			Code: global.ErrorNodeNotExists,
			Msg:  global.ErrorNodeNotExists.String(),
		}
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &nodeStatus); err != nil {
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  err.Error(),
		}
	}

	if nodeStatus.Online {
		return &global.Error{
			Code: global.ErrorNodeIsRunning,
			Msg:  global.ErrorNodeIsRunning.String(),
		}
	}

	//创建tasks cursor遍历节点上所有任务
	cursor, err := resourceutils.NewEtcdCursor(global.GetEtcdClient(), global.TasksNodePrefix+model.NodeID, 5)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  err.Error(),
		}
	}
	// 销毁节点上任务逻辑
	if strings.ToLower(model.TasksOnNodePolice) == "destroy" {
		//遍历节点上任务并删除
		for !cursor.IsFinished() {
			kvs, err := cursor.Next()
			if err != nil {
				global.RSPLog.Sugar().Error(err)
				continue
			}

			for _, v := range kvs {
				var tasksNodeVal global.TasksNodeVal
				if err := json.Unmarshal(v.Value, &tasksNodeVal); err != nil {
					global.RSPLog.Sugar().Error(err)
					continue
				}

				if err := RemoveTask(tasksNodeVal.TaskID); err != nil {
					global.RSPLog.Sugar().Error(err)
					continue
				}
			}
		}
	}

	// 迁移节点上任务逻辑
	if strings.ToLower(model.TasksOnNodePolice) == "migrate" {
		//节点选择器
		selector := node.NewSelector()
		for !cursor.IsFinished() {
			kvs, err := cursor.Next()
			if err != nil {
				global.RSPLog.Sugar().Error(err)
				continue
			}

			nodelist, err := selector.SelectNode()
			if err != nil {
				global.RSPLog.Sugar().Error(err)
				continue
			}

			var taskIDArray []string
			for _, v := range kvs {
				var tasksNodeVal global.TasksNodeVal
				if err := json.Unmarshal(v.Value, &tasksNodeVal); err != nil {
					global.RSPLog.Sugar().Error(err)
					continue
				}
				taskIDArray = append(taskIDArray, tasksNodeVal.TaskID)
			}
			TaskMigrate(taskIDArray, global.NodeTypeRedissyncer, (*nodelist)[0].Key)
		}
	}

	if strings.ToLower(model.TasksOnNodePolice) != "migrate" && strings.ToLower(model.TasksOnNodePolice) != "destroy" {
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  errors.New("must set tasksOnNodePolice").Error(),
		}
	}

	// 删除节点逻辑
	if _, err := global.GetEtcdClient().Delete(context.Background(), global.NodesPrefix+model.NodeType+"/"+model.NodeID); err != nil {
		global.RSPLog.Sugar().Error(err)
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  err.Error(),
		}
	}

	return nil

}
