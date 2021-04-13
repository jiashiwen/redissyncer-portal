// Etcd游标
// 根据查询前缀和返回PageSize 定义游标
// 游标单向向前
package resourceutils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	uuid "github.com/satori/go.uuid"
	"redissyncer-portal/commons"
	"redissyncer-portal/global"
	"redissyncer-portal/node"
	"sync"
)

var cursorQueryMap map[string]*EtcdCursor
var once sync.Once

type EtcdCursor struct {
	QueryID            string
	EtcdLeaseID        clientv3.LeaseID
	CurrentPage        int64
	EtcdPaginte        *EtcdPaginte
	LastQueryTimeStamp int64 //unix 时间戳 毫秒
	QueryFinish        bool
}

//获取queryIDMap
func GetCursorQueryMap() *map[string]*EtcdCursor {
	once.Do(func() {
		cursorQueryMap = make(map[string]*EtcdCursor)
	})
	return &cursorQueryMap
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
		cursor.LogoutFromCursorMap()
	}
	cursor.CurrentPage = cursor.EtcdPaginte.CurrentPage
	cursor.LastQueryTimeStamp = commons.GetCurrentUnixMillisecond()
	return kv, err
}

//注册到本地cursorMap
func (cursor *EtcdCursor) RegisterToCursorMap() error {
	cursorMap := *GetCursorQueryMap()
	cursorMap[cursor.QueryID] = cursor
	return nil
}

//从本地cursorMap注销
func (cursor *EtcdCursor) LogoutFromCursorMap() {
	cursorMap := *GetCursorQueryMap()
	delete(cursorMap, cursor.QueryID)
}

//通过queryID获得cursor指针
func GetCursorByQueryID(queryID string) (*EtcdCursor, error) {
	cursorMap := GetCursorQueryMap()
	cursor := (*cursorMap)[queryID]
	if cursor != nil {
		return cursor, nil
	}
	return nil, errors.New("cursor not exists")
}

//注册到etcd
func (cursor *EtcdCursor) RegisterToEtcd(cli *clientv3.Client) error {
	lease := clientv3.NewLease(cli)
	//未注册过的cursor生成新租约注册
	if cursor.EtcdLeaseID == 0 {
		gResp, err := lease.Grant(context.Background(), 300)
		if err != nil {
			return err
		}
		cursor.EtcdLeaseID = gResp.ID
		valJson, _ := json.Marshal(global.GetNodeInfo())
		if _, err := cli.Put(context.Background(), global.CursorPrefix+cursor.QueryID, string(valJson), clientv3.WithLease(cursor.EtcdLeaseID)); err != nil {
			return err
		}
	}

	//已注册过的cursor进行续约，若续约前已过期则把EtcdLeaseID置0，待下次注册重新生成租约
	if _, err := lease.KeepAliveOnce(context.Background(), cursor.EtcdLeaseID); err != nil {
		cursor.EtcdLeaseID = 0
		return err
	}
	return nil
}

//从etcd查询queryID所在节点
func GetCursorNode(cli *clientv3.Client, queryID string) (*node.NodeStatus, error) {
	var nodeStatus node.NodeStatus
	resp, err := cli.Get(context.Background(), global.CursorPrefix+queryID)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, errors.New("cursor not exists on any node")
	}
	if err := json.Unmarshal(resp.Kvs[0].Value, &nodeStatus); err != nil {
		return nil, err
	}
	return &nodeStatus, nil

}

//是否已完成查询
func (cursor *EtcdCursor) IsFinished() bool {
	return cursor.QueryFinish
}
