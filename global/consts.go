package global

// TaskStatusType 任务状态类型
type TaskStatusType int

// TaskType 任务类型
type TaskType int

const (
	TasksTaskIDPrefix     = "/tasks/taskid/"     //TasksTaskidPrefix 任务id前缀 key:/tasks/taskid/{taskid} ;value: taskstatusjson
	TasksNodePrefix       = "/tasks/node/"       //TasksNodePrefix key:/tasks/node/{nodeId}/{taskId}; value:{"nodeId":"xxx","taskId":"xxx"}
	TasksGroupIDPrefix    = "/tasks/groupid/"    //TasksGroupidPrefix key:/tasks/groupid/{groupid}/{taskId};value:{"groupId":"xxx","taskId":"xxx"}
	TasksStatusPrefix     = "/tasks/status/"     //TasksStatusPrefix key:/tasks/status/{currentstatus}/{taskid};value:{"taskId":"testId"}
	TasksRdbversionPrefix = "/tasks/rdbversion/" //TasksRdbversionPrefix key:/tasks/rdbversion/{redisVersion}/{rdbVersion};value:{"id":1,"redis_version": "2.6","rdb_version": 6}
	TasksOffsetPrefix     = "/tasks/offset/"     //TasksOffsetPrefix key:/tasks/offset/{taskId};value:{"replId":"xxx","replOffset":"-1"}
	TasksNamePrefix       = "/tasks/name/"       //TasksNamePrefix key:/tasks/name/{taskname};value:{"taskId":"testId"}
	TasksTypePrefix       = "/tasks/type/"       //TasksTypePrefix key:/tasks/type/{type}/{taskId};value:{"taskid":"xxx","groupId":"xxx","nodeId":"xxx"}
	TasksBigkeyPrefix     = "/tasks/bigkey/"     //TasksBigkeyPrefix key:/tasks/bigkey/{taskId}/{bigKey};value:{"id":1,"taskId":"xxx","command":"xxx","command_type":"xxx"}
	TasksMd5Prefix        = "/tasks/md5/"        // TasksMd5Prefix key:/tasks/md5/{md5};value:{"taskid":"xxx","groupId":"xxx","nodeId":"xx"}
	NodesPrefix           = "/nodes/"            // NodesPrefix key:/nodes/{nodetype}/{nodeID};value:{"nodeaddr":"10.0.0.1","nodeport":8888,"online":true,"lastreporttime":1615275888064}
	CursorPrefix          = "/cursor/"           // CursorPrefix key:/cursor/{qureryID};value:{"nodeaddr":"10.0.0.1","nodeport":8888,"online":true,"lastreporttime":1615275888064}

)

const (
	//IDSeed id种子
	IDSeed = "/uniqid/idseed"
)

const (
	// NodeTypePortal 节点类型 portal
	NodeTypePortal = "portal"
	// NodeTypeRedissyncer 节点类型 redissyncernodeserver
	NodeTypeRedissyncer = "redissyncernodeserver"
)

//任务状态
const (
	// TaskStatusTypeSTOP STOP	0	任务停止状态
	TaskStatusTypeSTOP TaskStatusType = 0
	// TaskStatusTypeCREATING CREATING	1	任务创建中
	TaskStatusTypeCREATING TaskStatusType = 1
	// TaskStatusTypeCREATED CREATED	2	任务创建完成
	TaskStatusTypeCREATED TaskStatusType = 2
	// TaskStatusTypeSTARTING 3	任务运行中，表示数据同步以前，发送psync命令，源redis进行bgsave 生成rdb的过程；描述不太贴切，待改进
	TaskStatusTypeSTARTING TaskStatusType = 3
	// TaskStatusTypeBROKEN BROKEN	5	任务异常
	TaskStatusTypeBROKEN TaskStatusType = 5
	// TaskStatusTypeRDBRUNNING RDBRUNNING	6	全量RDB同步过程中
	TaskStatusTypeRDBRUNNING TaskStatusType = 6
	// TaskStatusTypeCOMMANDRUNNING  COMMANDRUNNING	7	增量同步中
	TaskStatusTypeCOMMANDRUNNING TaskStatusType = 7
	// TaskStatusTypeFINISH FINISH	8
	TaskStatusTypeFINISH TaskStatusType = 8
)

//任务类型
const (
	// TaskTypeSYNC SYNC 1 replication 已使用
	TaskTypeSYNC TaskType = 1
	// TaskTypeRDB RDB 2 RDB文件解析 已使用
	TaskTypeRDB TaskType = 2
	// TaskTypeAOF AOF 3 AOF文件解析 已使用
	TaskTypeAOF TaskType = 3
	// TaskTypeMIXED MIXED 4 混合文件解析 已使用
	TaskTypeMIXED TaskType = 4
	// TaskTypeONLINERDB ONLINERDB 5 在线RDB解析 已使用
	TaskTypeONLINERDB TaskType = 5
	// TaskTypeONLINEAOF ONLINEAOF 6 在线AOF 已使用
	TaskTypeONLINEAOF TaskType = 6
	// TaskTypeONLINEMIXED ONLINEMIXED 7 在线混合文件解析 已使用
	TaskTypeONLINEMIXED TaskType = 7
	// TaskTypeCOMMANDDUMPUP COMMANDDUMPUP 8 增量命令实时备份 已使用
	TaskTypeCOMMANDDUMPUP TaskType = 8
)

func (taskStatusType TaskStatusType) String() string {
	switch taskStatusType {
	case TaskStatusTypeSTOP:
		return "STOP"
	case TaskStatusTypeCREATING:
		return "CREATING"
	case TaskStatusTypeCREATED:
		return "CREATED"
	case TaskStatusTypeSTARTING:
		return "STARTING"
	case TaskStatusTypeBROKEN:
		return "BROKEN"
	case TaskStatusTypeRDBRUNNING:
		return "RDBRUNNING"
	case TaskStatusTypeCOMMANDRUNNING:
		return "COMMANDRUNNING"
	case TaskStatusTypeFINISH:
		return "FINISH"
	default:
		return ""
	}
}

func (taskType TaskType) String() string {
	switch taskType {
	case TaskTypeSYNC:
		return "SYNC"
	case TaskTypeRDB:
		return "RDB"
	case TaskTypeAOF:
		return "AOF"
	case TaskTypeMIXED:
		return "MIXED"
	case TaskTypeONLINERDB:
		return "ONLINERDB"
	case TaskTypeONLINEAOF:
		return "ONLINEAOF"
	case TaskTypeONLINEMIXED:
		return "ONLINEMIXED"
	case TaskTypeCOMMANDDUMPUP:
		return "COMMANDDUMPUP"
	default:
		return ""
	}
}

// TaskStatus 任务状态
type TaskStatus struct {
	Afresh             bool        `mapstructure:"afresh" json:"afresh" yaml:"afresh"`
	AllKeyCount        int64       `mapstructure:"allKeyCount" json:"allKeyCount" yaml:"allKeyCount"`
	AutoStart          bool        `mapstructure:"autostart" json:"autostart" yaml:"autostart"`
	BatchSize          int64       `mapstructure:"batchSize" json:"batchSize" yaml:"batchSize"`
	CommandFilter      string      `mapstructure:"commandFilter" json:"commandFilter" yaml:"commandFilter"`
	CreateTime         string      `mapstructure:"createTime" json:"createTime" yaml:"createTime"`
	DBMapper           interface{} `mapstructure:"dbMapper" json:"dbMapper" yaml:"dbMapper"`
	DBMapping          interface{} `mapstructure:"dbMapping" json:"dbMapping" yaml:"dbMapping"`
	ErrorCount         int64       `mapstructure:"errorCount" json:"errorCount" yaml:"errorCount"`
	ExpandJSON         string      `mapstructure:"expandJson" json:"expandJson" yaml:"expandJson"`
	FileAddress        string      `mapstructure:"fileAddress" json:"fileAddress" yaml:"fileAddress"`
	FilterType         string      `mapstructure:"filterType" json:"filterType" yaml:"filterType"`
	GroupID            string      `mapstructure:"groupId" json:"groupId" yaml:"groupId"`
	ID                 string      `mapstructure:"id" json:"id" yaml:"id"`
	KeyFilter          string      `mapstructure:"keyFilter" json:"keyFilter" yaml:"keyFilter"`
	LastKeyCommitTime  int64       `mapstructure:"lastKeyCommitTime" json:"lastKeyCommitTime" yaml:"lastKeyCommitTime"`
	LastKeyUpdateTime  int64       `mapstructure:"lastKeyUpdateTime" json:"lastKeyUpdateTime" yaml:"lastKeyUpdateTime"`
	MD5                string      `mapstructure:"md5" json:"md5" yaml:"md5"`
	NodeID             string      `mapstructure:"nodeId" json:"nodeId" yaml:"nodeId"`
	Offset             int64       `mapstructure:"offset" json:"offset" yaml:"offset"`
	OffsetPlace        int         `mapstructure:"offsetPlace" json:"offsetPlace" yaml:"offsetPlace"`
	RdbKeyCount        int64       `mapstructure:"rdbKeyCount" json:"rdbKeyCount" yaml:"rdbKeyCount"`
	RdbVersion         int         `mapstructure:"rdbVersion" json:"rdbVersion" yaml:"rdbVersion"`
	RealKeyCount       int64       `mapstructure:"realKeyCount" json:"realKeyCount" yaml:"realKeyCount"`
	RedisVersion       float64     `mapstructure:"redisVersion" json:"redisVersion" yaml:"redisVersion"`
	ReplID             string      `mapstructure:"replId" json:"replId" yaml:"replId"`
	SourceACL          bool        `mapstructure:"sourceAcl" json:"sourceAcl" yaml:"sourceAcl"`
	SourceHost         string      `mapstructure:"sourceHost" json:"sourceHost" yaml:"sourceHost"`
	SourcePassword     string      `mapstructure:"sourcePassword" json:"sourcePassword" yaml:"sourcePassword"`
	SourcePort         int         `mapstructure:"sourcePort" json:"sourcePort" yaml:"sourcePort"`
	SourceRedisAddress string      `mapstructure:"sourceRedisAddress" json:"sourceRedisAddress" yaml:"sourceRedisAddress"`
	SourceRedisType    int         `mapstructure:"sourceRedisType" json:"sourceRedisType" yaml:"sourceRedisType"`
	SourceURI          string      `mapstructure:"sourceUri" json:"sourceUri" yaml:"sourceUri"`
	SourceUserName     string      `mapstructure:"sourceUserName" json:"sourceUserName" yaml:"sourceUserName"`
	Status             int         `mapstructure:"status" json:"status" yaml:"status"`
	SyncType           int         `mapstructure:"syncType" json:"syncType" yaml:"syncType"`
	TargetACL          bool        `mapstructure:"targetAcl" json:"targetAcl" yaml:"targetAcl"`
	TargetHost         string      `mapstructure:"targetHost" json:"targetHost" yaml:"targetHost"`
	TargetPassword     string      `mapstructure:"targetPassword" json:"targetPassword" yaml:"targetPassword"`
	TargetPort         int         `mapstructure:"targetPort" json:"targetPort" yaml:"targetPort"`
	TargetRedisAddress string      `mapstructure:"targetRedisAddress" json:"targetRedisAddress" yaml:"targetRedisAddress"`
	TargetRedisType    int         `mapstructure:"targetRedisType" json:"targetRedisType" yaml:"targetRedisType"`
	TargetURI          []string    `mapstructure:"targetUri" json:"targetUri" yaml:"targetUri"`
	TargetUserName     string      `mapstructure:"targetUserName" json:"targetUserName" yaml:"targetUserName"`
	TaskID             string      `mapstructure:"taskId" json:"taskId" yaml:"taskId"`
	TaskMsg            string      `mapstructure:"taskMsg" json:"taskMsg" yaml:"taskMsg"`
	TaskName           string      `mapstructure:"taskName" json:"taskName" yaml:"taskName"`
	TaskType           int         `mapstructure:"tasktype" json:"tasktype" yaml:"tasktype"`
	TimeDeviation      int64       `mapstructure:"timeDeviation" json:"timeDeviation" yaml:"timeDeviation"`
	UpdateTime         string      `mapstructure:"updateTime" json:"updateTime" yaml:"updateTime"`
}

type TasksOffset struct {
	ReplID     string `mapstructure:"replId" json:"replId" yaml:"replId"`
	ReplOffset int64  `mapstructure:"replOffset" json:"replOffset" yaml:"replOffset"`
}

type TaskIDVal struct {
	TaskID string `mapstructure:"taskId" json:"taskId" yaml:"taskId"`
}

type TasksNodeVal struct {
	NodeID string `mapstructure:"NodeID" json:"NodeID" yaml:"NodeID"`
	TaskID string `mapstructure:"taskId" json:"taskId" yaml:"taskId"`
}

type TasksGroupIDVal struct {
	GroupID string `mapstructure:"groupId" json:"groupId" yaml:"groupId"`
	TaskID  string `mapstructure:"taskId" json:"taskId" yaml:"taskId"`
}

type TasksTypeVal struct {
	TaskID  string `map:"taskId" json:"taskId" yaml:"taskId"`
	GroupID string `map:"groupId" json:"groupId" yaml:"groupId"`
	NodeID  string `map:"nodeId" json:"nodeId" yaml:"nodeId"`
}

type TasksMD5Val struct {
	TaskID  string `map:"taskId" json:"taskId" yaml:"taskId"`
	GroupID string `map:"groupId" json:"groupId" yaml:"groupId"`
	NodeID  string `map:"nodeId" json:"nodeId" yaml:"nodeId"`
}
