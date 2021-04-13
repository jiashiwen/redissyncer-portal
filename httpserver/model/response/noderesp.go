package response

import (
	"redissyncer-portal/global"
	"redissyncer-portal/node"
)

type NodeResult struct {
	NodeID     string           `mapstructure:"nodeID" json:"nodeID" yaml:"nodeID"`
	NodeType   string           `mapstructure:"nodeType" json:"nodeType" yaml:"nodeType"`
	Errors     *global.Error    `mapstructure:"errors" json:"errors" yaml:"errors"`
	NodeStatus *node.NodeStatus `mapstructure:"nodeStatus" json:"nodeStatus" yaml:"nodeStatus"`
}
