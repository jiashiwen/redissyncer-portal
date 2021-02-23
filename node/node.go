package node

import (
	"context"
	"encoding/json"
	"errors"
	"redissyncer-portal/global"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

const (
	NodesKeyPrefix = "/nodes/"
)

type Node struct {
	//节点类型
	NodeType string `mapstructure:"nodetype" json:"nodetype" yaml:"nodetype"`

	//节点Id
	NodeId string `mapstructure:"nodeid" json:"nodeid" yaml:"nodeid"`

	//节点ip地址
	NodeAddr string `mapstructure:"nodeaddr" json:"nodeaddr" yaml:"nodeaddr"`

	//探活port
	NodePort int `mapstructure:"nodeport" json:"nodeport" yaml:"nodeport"`

	//etcd 客户端
	EtcdClient *clientv3.Client

	//轮询Ticker
	NodeTicker *time.Ticker

	//Node context
	NodeContext context.Context

	//巡检器cancle
	NodeCancel context.CancelFunc
}

type NodeStatus struct {
	//节点ip地址
	NodeAddr string `mapstructure:"nodeaddr" json:"nodeaddr" yaml:"nodeaddr"`

	//探活port
	NodePort int `mapstructure:"nodeport" json:"nodeport" yaml:"nodeport"`

	//是否在线
	Online bool `mapstructure:"online" json:"online" yaml:"online"`

	//最后上报时间，unix时间戳
	LastReportTime int64 `mapstructure:"lastreporttime" json:"lastreporttime" yaml:"lastreporttime"`
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

	wg.Add(1)
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

	//若key存在，Online为true返回onde exists 错误
	if len(getResp.Kvs) > 0 {
		var nodeStatus NodeStatus
		if err := json.Unmarshal(getResp.Kvs[0].Value, &nodeStatus); err != nil {
			return err
		}

		if nodeStatus.Online == true {
			return errors.New("node exists")
		}

	}

	nodeStatus := &NodeStatus{
		NodeAddr:       node.NodeAddr,
		NodePort:       node.NodePort,
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
			NodeAddr:       node.NodeAddr,
			NodePort:       node.NodePort,
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
