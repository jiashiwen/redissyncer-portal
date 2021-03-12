package service

import (
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"testing"
)

func TestGetTaskStatus(t *testing.T) {
	config := "../../config.yaml"
	global.RSPViper = core.Viper(config)
	global.RSPLog = core.Zap()
	ids := make([]string, 5)
	ids = append(ids, "8A3EB419098547258D94EE8BDFE49F3C")

	idsmap, err := GetTaskStatus(ids)

	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", idsmap)

}
