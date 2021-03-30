package service

import (
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"testing"
)

var config string = "../../config.yaml"

func TestNodeAllTypes(t *testing.T) {
	global.RSPViper = core.Viper(config)
	NodeAllTypes()
}
