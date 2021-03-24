package service

import (
	"context"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"strings"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

var config string = "../../config.yaml"

func TestRemoveTask(t *testing.T) {
	global.RSPViper = core.Viper(config)
	global.RSPLog = core.Zap()
	resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithLimit(3))

	if err != nil {
		t.Error(err)
		return
	}

	taskids := []string{}
	for _, v := range resp.Kvs {
		taskids = append(taskids, strings.Split(string(v.Key), "/")[3])
	}

	for _, v := range taskids {
		if err := RemoveTask(v); err != nil {
			t.Error(err)
		}
	}

}

func TestGetTaskStatusByName(t *testing.T) {

	global.RSPViper = core.Viper(config)
	global.RSPLog = core.Zap()
	resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksNamePrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithLimit(3))

	if err != nil {
		t.Error(err)
		return
	}
	names := []string{}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			names = append(names, strings.Split(string(v.Key), "/")[3])
		}
	}

	t.Log(names)

	taskstatusArry, err := GetTaskStatusByName(names)

	if err != nil {
		t.Error(err)
	}

	for _, v := range taskstatusArry {
		t.Logf("%+v", v)
	}

}

func TestGetTaskStatusByGroupID(t *testing.T) {

	global.RSPViper = core.Viper(config)
	global.RSPLog = core.Zap()
	resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksGroupIDPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithLimit(3))

	if err != nil {
		t.Error(err)
		return
	}
	groupids := []string{}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			groupids = append(groupids, strings.Split(string(v.Key), "/")[3])
		}
	}

	t.Log(groupids)

	taskstatusArry, err := GetTaskStatusByGroupIDs(groupids)

	if err != nil {
		t.Error(err)
	}

	for _, v := range taskstatusArry {
		t.Logf("%+v", v)
	}

}

func TestGetTaskStatus(t *testing.T) {
	// config := "../../config.yaml"
	global.RSPViper = core.Viper(config)
	global.RSPLog = core.Zap()
	resp, err := global.GetEtcdClient().Get(context.Background(), global.TasksTaskIDPrefix, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithLimit(3))

	if err != nil {
		t.Error(err)
		return
	}
	ids := []string{}
	if len(resp.Kvs) > 0 {
		for _, v := range resp.Kvs {
			ids = append(ids, strings.Split(string(v.Key), "/")[3])
		}
	}

	t.Log(ids)

	taskstatusArry, err := GetTaskStatusByIDs(ids)

	if err != nil {
		t.Error(err)
	}

	for _, v := range taskstatusArry {
		t.Logf("%+v", v)
	}

}
