package config

type Server struct {
	//JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap  Zap  `mapstructure:"zap" json:"zap" yaml:"zap"`
	Etcd Etcd `mapstructure:"etcd" json:"etcd" yaml:"etcd"`
	//Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	//Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	//Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	//System  System  `mapstructure:"system" json:"system" yaml:"system"`
	//Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	//Local Local `mapstructure:"local" json:"local" yaml:"local"`
}
