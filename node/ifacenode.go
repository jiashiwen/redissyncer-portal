package node

//node 用于启动node进程注册节点并定时上报节点最新时间戳

type inode interface {
	//启动节点
	Start()
	//停止节点
	Stop()
	//注册节点
	Registry()
	//上报节点状态
	ReportStatus()
}
