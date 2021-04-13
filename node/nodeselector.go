// 节点选择器，用于在创建或迁移任务时选择资源占用最少的节点
package node

import (
	"context"
	"encoding/json"
	"redissyncer-portal/commons"
	"redissyncer-portal/global"
	"strings"

	"github.com/coreos/etcd/clientv3"
)

type Selector struct {
	EtcdClient *clientv3.Client
}

func NewSelector() *Selector {
	return &Selector{global.GetEtcdClient()}
}

//SelectNode 节点选择器，根据规则选择合适的节点来承载业务
//简单做法是选择任务数量最少的节点
//返回值为节点id及其任务数量的列表
func (nodeSelector *Selector) SelectNode() (*commons.PairList, error) {
	nodeIncludeTasks := make(map[string]int64)

	//解决节点启动后没有任务的问题
	//将所有worker节点预制到map，值为0
	//获取所有类型为 redissyncer 的存活node节点
	nodeResp, err := nodeSelector.EtcdClient.Get(context.Background(), global.NodesPrefix+global.NodeTypeRedissyncer, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range nodeResp.Kvs {
		var nodeStatus NodeStatus
		if err := json.Unmarshal(v.Value, &nodeStatus); err != nil {
			global.RSPLog.Sugar().Error(err)
			continue
		}
		if nodeStatus.Online {
			nodeIncludeTasks[nodeStatus.NodeID] = 0
		}
	}

	//获取node节点所有任务数量的map[nodeid]数量
	getResp, err := nodeSelector.EtcdClient.Get(context.Background(), global.TasksNodePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range getResp.Kvs {
		nodeID := strings.Split(string(v.Key), "/")[3]
		if val, ok := nodeIncludeTasks[nodeID]; ok {
			nodeIncludeTasks[nodeID] = val + 1
		}
	}

	//根据任务数量进行排序
	pairList := commons.SortMapByValue(nodeIncludeTasks, false)

	//返回列表结果
	return &pairList, nil

}
