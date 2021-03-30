package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	set "github.com/deckarep/golang-set"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver/model"
	"redissyncer-portal/httpserver/model/response"
	"redissyncer-portal/node"
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
		return &global.Error{
			Code: global.ErrorSystemError,
			Msg:  global.ErrorSystemError.String(),
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
			Msg:  global.ErrorSystemError.String(),
		}
	}

	if nodeStatus.Online {
		return &global.Error{
			Code: global.ErrorNodeIsRunning,
			Msg:  global.ErrorNodeIsRunning.String(),
		}
	}

	//Todo 销毁节点上任务逻辑
	if strings.ToLower(model.TasksOnNodePolice) == "destroy" {
		//遍历节点上任务并删除
	}

	//ToDo 迁移节点上任务逻辑
	if strings.ToLower(model.TasksOnNodePolice) == "migrate" {
		//节点选择
		//节点上任务遍历并迁移至合规节点
	}

	//Todo 删除节点逻辑

	return nil

}
