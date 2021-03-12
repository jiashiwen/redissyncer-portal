// 节点选择器，用于在创建或迁移任务时选择资源占用最少的节点
package node

import (
	"context"
	"redissyncer-portal/commons"
	"redissyncer-portal/global"
	"strings"

	"github.com/coreos/etcd/clientv3"
)

type NodeSelector struct {
	EtcdClient *clientv3.Client
}

//SelectNode 节点选择器，根据规则选择合适的节点来承载业务
//简单做法是选择任务数量最少的节点
//返回值为节点id及其任务数量的列表
func (nodeSelector *NodeSelector) SelectNode() (*commons.PairList, error) {

	//获取node节点所有任务数量的map[nodeid]数量
	nodeIncludeTasks := make(map[string]int64)
	getResp, err := nodeSelector.EtcdClient.Get(context.Background(), global.TasksNodePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range getResp.Kvs {
		nodeID := strings.Split(string(v.Key), "/")[3]
		nodeIncludeTasks[nodeID] = nodeIncludeTasks[nodeID] + 1
	}

	//根据任务数量进行排序
	pairList := commons.SortMapByValue(nodeIncludeTasks, false)

	//返回列表结果
	return &pairList, nil

}
