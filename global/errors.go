package global

type ErrorCode int

const (
	ErrorSystemError        ErrorCode = 10001
	ErrorNodeNotExists      ErrorCode = 40001
	ErrorTaskNotExists      ErrorCode = 50001
	ErrorTaskGroupNotExists ErrorCode = 50002
)

func (code ErrorCode) String() string {
	switch code {
	case 40001:
		return "node not exists"
	case 50001:
		return "task not exists"
	case 50002:
		return "task group not exists"
	default:
		return ""
	}
}
