//用于访问redissyncer-server
//任务操作包括：创建任务、启动任务、任务查询、任务停止和任务删除

package httpquerry

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	UrlLogin       = "/login"
	UrlCreateTask  = "/api/v2/createtask"
	UrlStartTask   = "/api/v2/starttask"
	UrlStopTask    = "/api/v2/stoptask"
	UrlRemoveTask  = "/api/v2/removetask"
	UrlListTasks   = "/api/v2/listtasks"
	ImportFilePath = "/api/v2/file/createtask"
)

type HttpRequest struct {
	Server     string
	Api        string
	Body       string
	HttpClient *http.Client
}

func New(server string) *HttpRequest {
	httpRequest := &HttpRequest{
		Server:     server,
		HttpClient: &http.Client{},
	}
	httpRequest.HttpClient.Timeout = 15 * time.Second
	return httpRequest
}

func (r *HttpRequest) ExecRequest() (string, error) {

	req, err := http.NewRequest("POST", r.Server+r.Api, strings.NewReader(r.Body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token", viper.GetString("token"))

	resp, respErr := r.HttpClient.Do(req)

	if respErr != nil {
		return "", respErr
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	bodyStr, jsonErr := json.MarshalIndent(dat, "", " ")
	if jsonErr != nil {
		return "", jsonErr
	}
	return string(bodyStr), nil
}

//登录
func (r *HttpRequest) Login(username, password string) (string, error) {
	jsonMap := make(map[string]interface{})
	jsonMap["username"] = username
	jsonMap["password"] = password
	loginJson, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	r.Api = UrlLogin
	r.Body = string(loginJson)
	return r.ExecRequest()
}

//创建同步任务
func (r *HttpRequest) CreateTask(createJson string) ([]string, error) {

	r.Api = UrlCreateTask
	r.Body = createJson

	resp, err := r.ExecRequest()
	if err != nil {
		return nil, err
	}
	taskIds := gjson.Get(resp, "data").Array()
	if len(taskIds) == 0 {
		return nil, errors.New("task create faile")
	}
	var taskIdsStrArray []string
	for _, v := range taskIds {
		taskIdsStrArray = append(taskIdsStrArray, gjson.Get(v.String(), "taskId").String())
	}

	return taskIdsStrArray, nil

}

//Start task
func (r *HttpRequest) StartTask(taskid string) (string, error) {
	jsonMap := make(map[string]interface{})
	jsonMap["taskid"] = taskid
	startJson, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	r.Api = UrlStartTask
	r.Body = string(startJson)
	return r.ExecRequest()

}

//Stop task by task ids
func (r *HttpRequest) StopTaskByIds(ids []string) (string, error) {
	jsonMap := make(map[string]interface{})
	jsonMap["taskids"] = ids
	stopJsonStr, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	r.Api = UrlStopTask
	r.Body = string(stopJsonStr)
	return r.ExecRequest()

}

//Remove task by name
func (r *HttpRequest) RemoveTaskByName(taskName string) (string, error) {
	jsonMap := make(map[string]interface{})

	taskids, err := r.GetSameTaskNameIds(taskName)
	if err != nil {
		return "", err
	}

	if len(taskids) == 0 {
		return "", errors.New("no taskid")
	}

	jsonMap["taskids"] = taskids
	stopJsonStr, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	r.Api = UrlStopTask
	r.Body = string(stopJsonStr)
	r.ExecRequest()

	r.Api = UrlRemoveTask
	r.Body = string(stopJsonStr)

	return r.ExecRequest()

}

//获取同步任务状态
func (r *HttpRequest) GetTaskStatus(ids []string) (map[string]string, error) {
	jsonMap := make(map[string]interface{})

	jsonMap["regulation"] = "byids"
	jsonMap["taskids"] = ids

	listTaskJsonStr, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}

	r.Api = UrlListTasks
	r.Body = string(listTaskJsonStr)

	listResp, err := r.ExecRequest()
	taskArray := gjson.Get(listResp, "data").Array()

	if len(taskArray) == 0 {
		return nil, errors.New("No status return")
	}

	statusMap := make(map[string]string)

	for _, v := range taskArray {
		id := gjson.Get(v.String(), "taskId").String()
		status := gjson.Get(v.String(), "status").String()
		statusMap[id] = status
	}

	return statusMap, nil
}

// @title    GetSameTaskNameIds
// @description   获取同名任务列表
// @auth      Jsw             2020/7/1   10:57
// @param    taskName        string         "任务名称"
// @return    taskIds        []string         "任务id数组"
func (r *HttpRequest) GetSameTaskNameIds(taskName string) ([]string, error) {

	var existTaskIds []string
	listJsonMap := make(map[string]interface{})
	listJsonMap["regulation"] = "bynames"
	listJsonMap["tasknames"] = strings.Split(taskName, ",")
	listJsonStr, err := json.Marshal(listJsonMap)
	if err != nil {
		return nil, err
	}

	r.Api = UrlListTasks
	r.Body = string(listJsonStr)

	listResp, err := r.ExecRequest()
	if err != nil {
		return nil, err
	}
	taskList := gjson.Get(listResp, "data").Array()

	if len(taskList) > 0 {
		for _, v := range taskList {
			existTaskIds = append(existTaskIds, gjson.Get(v.String(), "taskId").String())
		}
	}
	return existTaskIds, nil
}
