package utils

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type EtcdClient struct {
	Client *clientv3.Client
	Config clientv3.Config
}

func NewEtcdClient(config clientv3.Config) (*EtcdClient, error) {

	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	etcdClient := &EtcdClient{
		Client: cli,
		Config: config,
	}
	return etcdClient, nil
}

func (eCli *EtcdClient) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := eCli.Client.Cluster.MemberList(ctx)
	cancel()
	return err
}
