package config

type NodeInfo struct {
	NodeType       string `mapstructure:"nodetype" json:"nodetype" yaml:"nodetype"`
	NodeId         string `mapstructure:"nodeid" json:"nodeid" yaml:"nodeid"`
	NodeAddr       string `mapstructure:"nodeaddr" json:"nodeaddr" yaml:"nodeaddr"`
	NodePort       int    `mapstructure:"nodeport" json:"nodeport" yaml:"nodeport"`
	NodeTickerTime int    `mapstructure:"nodetickertime" json:"nodetickertime" yaml:"nodeticertime"`
}
