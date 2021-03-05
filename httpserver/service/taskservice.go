package service

import (
	"redissyncer-portal/global"
	"redissyncer-portal/node"
)

func TaskCreate() error {
	selector := node.NodeSelector{
		EtcdClient: global.GetEtcdClient(),
	}

	pairelist, err := selector.SelectNode()

	if err != nil {
		return err
	}

	global.RSPLog.Sugar().Debug(pairelist)
	return nil

	//探活

	//发送创建请求

}
