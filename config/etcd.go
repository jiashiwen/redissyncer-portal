package config

import "time"

type Etcd struct {
	Endpoints   []string      `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"`
	DialTimeout time.Duration `mapstructure:"dialtimeout" json:"dialtimeout" yaml:"dialtimeout"`
	Username    string        `mapstructure:"username" json:"username" yaml:"username"`
	Password    string        `mapstructure:"password" json:"password" yaml:"password"`
}
