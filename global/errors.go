package global

type ErrorCode int
type Error struct {
	Code ErrorCode `map:"code" json:"code" yaml:"code"`
	Msg  string    `map:"msg" json:"msg" yaml:"msg"`
}

const (
	ErrorSystemError        ErrorCode = 10001
	ErrorCursorFinished     ErrorCode = 20001
	ErrorNodeNotExists      ErrorCode = 40001
	ErrorNodeIsRunning      ErrorCode = 40002
	ErrorNodeNotAlive       ErrorCode = 40003
	ErrorTaskNotExists      ErrorCode = 50001
	ErrorTaskStatusIsNil    ErrorCode = 50002
	ErrorTaskGroupNotExists ErrorCode = 50003
	ErrorEtcdKeyNotExists   ErrorCode = 60001
)

func (code ErrorCode) String() string {
	switch code {
	case 20001:
		return "cursor query have finished"
	case 40001:
		return "node not exists"
	case 40002:
		return "node is running"
	case 40003:
		return "node not alive"
	case 50001:
		return "task not exists"
	case 50002:
		return "task status is nil"
	case 50003:
		return "task group not exists"
	case 60001:
		return "etcd key not exists"
	default:
		return ""
	}
}

func (err *Error) Error() string {
	return err.Code.String()
}
