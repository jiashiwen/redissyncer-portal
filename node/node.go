package node

import (
	"context"
	"encoding/json"
	"errors"
	"redissyncer-portal/global"
	"redissyncer-portal/httpquerry"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	NodesKeyPrefix = "/nodes/"
)

//本节点描述
type Node struct {
	//节点类型
	NodeType string `map:"nodetype" json:"nodetype" yaml:"nodetype"`

	//节点Id
	NodeId string `map:"nodeid" json:"nodeid" yaml:"nodeid"`

	//节点ip地址
	NodeAddr string `map:"nodeaddr" json:"nodeaddr" yaml:"nodeaddr"`

	//探活port
	NodePort int `map:"nodeport" json:"nodeport" yaml:"nodeport"`

	//etcd 客户端
	EtcdClient *clientv3.Client

	//轮询Ticker
	NodeTicker *time.Ticker

	//Node context
	NodeContext context.Context

	//巡检器cancle
	NodeCancel context.CancelFunc
}

//本节点状态
type NodeStatus struct {
	//节点类型
	NodeType string `map:"nodetype" json:"nodetype" yaml:"nodetype"`

	//节点Id
	NodeID string `maps:"nodeid" json:"nodeid" yaml:"nodeid"`

	//节点ip地址
	NodeAddr string `maps:"nodeaddr" json:"nodeaddr" yaml:"nodeaddr"`

	//探活port
	NodePort int `map:"nodeport" json:"nodeport" yaml:"nodeport"`

	//探活url
	HeartbeatUrl string `map:"heartbeaturl" json:"heartbeaturl" yaml:"heartbeaturl"`

	//是否在线
	Online bool `map:"online" json:"online" yaml:"online"`

	//最后上报时间，unix时间戳
	LastReportTime int64 `map:"lastreporttime" json:"lastreporttime" yaml:"lastreporttime"`
}

//初始化node
func NewNode() *Node {
	ctx, cancel := context.WithCancel(context.Background())
	return &Node{
		NodeType:    global.GetNodeInfo().NodeType,
		NodeId:      global.GetNodeInfo().NodeId,
		NodeAddr:    global.GetNodeInfo().NodeAddr,
		NodePort:    global.GetNodeInfo().NodePort,
		EtcdClient:  global.GetEtcdClient(),
		NodeTicker:  time.NewTicker(time.Duration(global.GetNodeInfo().NodeTickerTime) * time.Second),
		NodeContext: ctx,
		NodeCancel:  cancel,
	}
}

func (node *Node) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	defer node.NodeCancel()
	node.ReportStatus(node.NodeContext, wg)

}

//停止节点
func (node *Node) Stop() {
	node.NodeCancel()
}

//注册节点
func (node *Node) Registry() error {
	nodeKey := NodesKeyPrefix + node.NodeType + "/" + node.NodeId
	getResp, err := node.EtcdClient.Get(context.TODO(), nodeKey)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}

	//若key存在，Online为true返回node exists 错误
	if len(getResp.Kvs) > 0 {
		var nodeStatus NodeStatus
		if err := json.Unmarshal(getResp.Kvs[0].Value, &nodeStatus); err != nil {
			return err
		}

		if nodeStatus.Online == true && httpquerry.NodeAlive(nodeStatus.NodeAddr, strconv.Itoa(nodeStatus.NodePort)) {
			return errors.New("node exists")
		}

	}

	nodeStatus := &NodeStatus{
		NodeType:       node.NodeType,
		NodeID:         node.NodeId,
		NodeAddr:       node.NodeAddr,
		NodePort:       node.NodePort,
		HeartbeatUrl:   "/health",
		Online:         true,
		LastReportTime: time.Now().UnixNano() / 1e6,
	}

	statusJson, err := json.Marshal(nodeStatus)
	if err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}

	if _, err := global.GetEtcdClient().Put(context.Background(), nodeKey, string(statusJson)); err != nil {
		global.RSPLog.Sugar().Error(err)
		return err
	}

	return nil

}

//上报节点状态
func (node *Node) ReportStatus(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	nodeKey := NodesKeyPrefix + node.NodeType + "/" + node.NodeId

	for range node.NodeTicker.C {
		now := time.Now().UnixNano() / 1e6
		nodeStatus := &NodeStatus{
			NodeType:       node.NodeType,
			NodeID:         node.NodeId,
			NodeAddr:       node.NodeAddr,
			NodePort:       node.NodePort,
			HeartbeatUrl:   "/health",
			Online:         true,
			LastReportTime: now,
		}
		select {
		case <-ctx.Done():
			return
		default:
			statusJson, err := json.Marshal(nodeStatus)

			if err != nil {
				global.RSPLog.Sugar().Error(err)
			}

			if _, err := global.GetEtcdClient().Put(context.Background(), nodeKey, string(statusJson)); err != nil {
				global.RSPLog.Sugar().Error(err)
			}
		}
	}
}
