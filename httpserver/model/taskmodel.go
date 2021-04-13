package model

type TaskCreateBody struct {
	Name     string `json:"name"`
	Password int64  `json:"password"`
}

type TaskStart struct {
	TaskID string `mapstructure:"taskid" json:"taskid" yaml:"taskid"`
	Afresh string `mapstructure:"afresh" json:"afresh" yaml:"afresh"`
}

type TaskStopBodyToNode struct {
	TaskIDs []string `maps:"taskids" json:"taskids" yaml:"taskids"`
}

type TaskIDBody struct {
	TaskID string `maps:"taskid" json:"taskid" yaml:"taskid"`
}

type TaskListByTaskIDs struct {
	TaskIDs []string `maps:"taskIDs" json:"taskIDs" yaml:"taskIDs"`
}

type TaskListByGroupIDs struct {
	GroupIDs []string `maps:"groupIDs" json:"groupIDs" yaml:"groupIDs"`
}

type TaskListByTaskNames struct {
	TaskNames []string `maps:"taskNames" json:"taskNames" yaml:"taskNames"`
}

type TaskListAll struct {
	QueryID   string `maps:"queryID" json:"queryID" yaml:"queryID"`
	BatchSize int64  `maps:"batchSize" json:"batchSize" yaml:"batchSize"`
	KeyPrefix string `maps:"keyPrefix" json:"keyPrefix" yaml:"keyPrefix"`
}

type TaskListByNode struct {
	NodeID    string `maps:"nodeID" json:"nodeID" yaml:"nodeID"`
	QueryID   string `maps:"queryID" json:"queryID" yaml:"queryID"`
	BatchSize int64  `maps:"batchSize" json:"batchSize" yaml:"batchSize"`
}
