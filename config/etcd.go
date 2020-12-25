package config

import "time"

type Etcd struct {
	Endpoints   []string      `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
	DialTimeout time.Duration `mapstructure:"dialtimeout" json:"dialtimeout" yaml:"dialtimeout"`
}
