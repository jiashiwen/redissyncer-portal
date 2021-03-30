package response

import (
	"redissyncer-portal/global"
	"redissyncer-portal/node"
)

type NodeResult struct {
	NodeID     string           `map:"nodeID" json:"nodeID" yaml:"nodeID"`
	NodeType   string           `map:"nodeType" json:"nodeType" yaml:"nodeType"`
	Errors     *global.Error    `map:"errors" json:"errors" yaml:"errors"`
	NodeStatus *node.NodeStatus `map:"nodeStatus" json:"nodeStatus" yaml:"nodeStatus"`
}
