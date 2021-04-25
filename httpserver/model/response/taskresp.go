package response

import "redissyncer-portal/global"

type StopTasksResult struct {
	TaskID string `map:"taskId" json:"taskId" yaml:"taskId"`
	Result string `map:"result" json:"result" yaml:"result"`
}

type TaskStatusResult struct {
	TaskID     string             `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors     *global.Error      `map:"errors" json:"errors" yaml:"errors"`
	TaskStatus *global.TaskStatus `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}

type TaskStatusResultByName struct {
	TaskName   string             `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors     *global.Error      `map:"errors" json:"errors" yaml:"errors"`
	TaskStatus *global.TaskStatus `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}

type TaskStatusResultByGroupID struct {
	GroupID         string              `map:"taskId" json:"taskId" yaml:"taskId"`
	Errors          *global.Error       `map:"errors" json:"errors" yaml:"errors"`
	TaskStatusArray []*TaskStatusResult `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}

type AllTaskStatusResult struct {
	QueryID         string              `map:"queryID" json:"queryID" yaml:"queryID"`
	LastPage        bool                `map:"lastPage" json:"lastPage" yaml:"lastPage"`
	CurrentPage     int64               `map:"currentPage" json:"currentPage" yaml:"currentPage"`
	Errors          []*global.Error     `map:"errors" json:"errors" yaml:"errors"`
	TaskStatusArray []*TaskStatusResult `map:"taskStatus" json:"taskStatus" yaml:"taskStatus"`
}
