package response

import "redissyncer-portal/global"

type StopTasksResult struct {
	TaskID string `map:"taskId" json:"taskId" yaml:"taskId"`
	Result string `map:"result" json:"result" yaml:"result"`
}

type TaskStatusResult struct {
	TaskID     string             `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors     *ErrorResult       `map:"errors" json:"errors" yaml:"errors"`
	TaskStatus *global.TaskStatus `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}

type TaskStatusResultByName struct {
	TaskName   string             `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors     *ErrorResult       `map:"errors" json:"errors" yaml:"errors"`
	TaskStatus *global.TaskStatus `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}

type TaskStatusResultByGroupID struct {
	GroupID         string              `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors          *ErrorResult        `map:"errors" json:"errors" yaml:"errors"`
	TaskStatusArray []*TaskStatusResult `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}
