// Etcd游标
// 根据查询前缀和返回pagesize 定义游标
// 游标单向向前
package resourceutils

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	uuid "github.com/satori/go.uuid"
	"redissyncer-portal/commons"
)

type EtcdCursor struct {
	QueryID            string
	EtcdPaginte        *EtcdPaginte
	LastQueryTimeStamp int64
	QueryFinish        bool
}

//新建EtcdCursor
func NewEtcdCursor(cli *clientv3.Client, keyPrefix string, pageSize int64) (*EtcdCursor, error) {
	etcdPaginte, err := NewEtcdPaginte(cli, keyPrefix, pageSize)
	if err != nil {
		return nil, err
	}

	etcdCursor := EtcdCursor{
		QueryID:            uuid.NewV4().String(),
		EtcdPaginte:        etcdPaginte,
		LastQueryTimeStamp: commons.GetCurrentUnixMillisecond(),
		QueryFinish:        false,
	}
	return &etcdCursor, nil
}

func (cursor *EtcdCursor) Next() ([]*mvccpb.KeyValue, error) {
	kv, err := cursor.EtcdPaginte.Next()
	if cursor.EtcdPaginte.LastPage {
		cursor.QueryFinish = true
	}
	cursor.LastQueryTimeStamp = commons.GetCurrentUnixMillisecond()
	return kv, err
}

//是否已完成查询
func (cursor *EtcdCursor) Finish() bool {
	return cursor.QueryFinish
}
