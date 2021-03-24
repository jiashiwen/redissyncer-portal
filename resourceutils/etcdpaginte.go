package resourceutils

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/pkg/errors"
	"redissyncer-portal/logger"
)

type EtcdPaginte struct {
	Kv             clientv3.KV
	Ctx            context.Context
	KeyPrefix      string
	TotalRecord    int64
	PageSize       int64
	Pages          int64
	Remainder      int64
	CurrentPage    int64
	FirstKeyArray  []string //每页首个key的数组
	CurrentLastKey string
	LastPage       bool
}

//初始化EtcdPaginte
func NewEtcdPaginte(cli *clientv3.Client, keyPrefix string, pageSize int64) (*EtcdPaginte, error) {

	if pageSize <= 0 {
		return nil, errors.New("PageSize must greater than 0")
	}

	ep := &EtcdPaginte{
		Kv:        clientv3.NewKV(cli),
		Ctx:       context.TODO(),
		KeyPrefix: keyPrefix,
		PageSize:  pageSize,
	}

	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithCountOnly(),
	}

	gr, err := ep.Kv.Get(context.TODO(), keyPrefix, opts...)

	if err != nil {
		return nil, err
	}

	ep.TotalRecord = gr.Count
	pages := ep.TotalRecord / pageSize
	remainder := ep.TotalRecord % pageSize

	if remainder > 0 {
		pages = pages + 1
	}

	if ep.TotalRecord == 0 {
		ep.LastPage = true
	} else {
		ep.LastPage = false
	}

	ep.Pages = pages
	ep.Remainder = remainder
	ep.CurrentPage = 0
	ep.FirstKeyArray = []string{}

	return ep, nil

}

//初始化EtcdPaginte同时遍历数据
func NewEtcdPaginteWithTraverse(cli *clientv3.Client, keyPrefix string, pageSize int64) (*EtcdPaginte, error) {

	etcdPaginte, err := NewEtcdPaginte(cli, keyPrefix, pageSize)

	if err != nil {
		return nil, err
	}

	if err := etcdPaginte.Traverse(); err != nil {
		return nil, err
	}

	return etcdPaginte, nil

}

func (ep *EtcdPaginte) Next() ([]*mvccpb.KeyValue, error) {

	if ep.TotalRecord == 0 {
		return nil, errors.New("no data returned")
	}

	if ep.LastPage {
		return nil, errors.New("last page,no next page")
	}

	//获取第一页数据的opts设置
	if ep.CurrentPage == 0 {
		opts := []clientv3.OpOption{
			clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		}

		//当总页数大于1页时，首页需要追加页限制条件
		if ep.Pages > 1 {
			opts = append(opts, clientv3.WithLimit(ep.PageSize))
		}

		gr, err := ep.Kv.Get(ep.Ctx, ep.KeyPrefix, opts...)
		if err != nil {
			return nil, err
		}

		//ep.CurrentFirstKey = string(gr.Kvs[0].Key)

		ep.CurrentLastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
		ep.CurrentPage = ep.CurrentPage + 1
		if int64(len(ep.FirstKeyArray)) < ep.CurrentPage {
			ep.FirstKeyArray = append(ep.FirstKeyArray, string(gr.Kvs[0].Key))
		} else {
			ep.FirstKeyArray[ep.CurrentPage-1] = string(gr.Kvs[0].Key)
		}
		//若返回数据仅有一页，添加末页标识为true
		if ep.Pages == 1 {
			ep.LastPage = true
		}
		return gr.Kvs, nil
	}

	//获取最后一页数据时，若最后一页记录不足pagesize则更改限制条件为余数+1
	limit := ep.PageSize + 1

	if ep.Remainder > 0 && ep.CurrentPage == ep.Pages-1 {
		limit = ep.Remainder + 1
	}

	opts := []clientv3.OpOption{
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(limit),
		clientv3.WithFromKey(),
	}
	gr, err := ep.Kv.Get(ep.Ctx, ep.CurrentLastKey, opts...)
	if err != nil {
		return nil, err
	}
	//ep.CurrentFirstKey = string(gr.Kvs[0].Key)

	ep.CurrentLastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
	ep.CurrentPage = ep.CurrentPage + 1
	if int64(len(ep.FirstKeyArray)) < ep.CurrentPage {
		ep.FirstKeyArray = append(ep.FirstKeyArray, string(gr.Kvs[1].Key))
	} else {
		ep.FirstKeyArray[ep.CurrentPage-1] = string(gr.Kvs[1].Key)
	}

	//若已是最后一页则更改末页标识为true
	if ep.CurrentPage == ep.Pages {
		ep.LastPage = true
	}

	return gr.Kvs[1:], nil

}

//遍历数据，填充ep.FirstKeyArray
func (ep *EtcdPaginte) Traverse() error {
	for {
		if ep.LastPage {
			return nil
		}

		if _, err := ep.Next(); err != nil {
			return err
		}
	}
}

//按页号获取数据
func (ep *EtcdPaginte) GetPage(page int64) ([]*mvccpb.KeyValue, error) {
	if page < 1 {
		return nil, errors.New("page must grater then 0")
	}

	if page > ep.Pages {
		return nil, errors.New("page must less than pages")
	}

	limit := ep.PageSize

	if page == ep.PageSize && ep.Remainder > 0 {
		limit = ep.Remainder
	}

	opts := []clientv3.OpOption{
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(limit),
		clientv3.WithFromKey(),
	}

	logger.Logger().Sugar().Info(ep.CurrentPage, "|", ep.FirstKeyArray[ep.CurrentPage-1])

	gr, err := ep.Kv.Get(ep.Ctx, ep.FirstKeyArray[page-1], opts...)
	if err != nil {
		return nil, err
	}

	ep.CurrentLastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
	ep.CurrentPage = page

	//若已是最后一页则更改末页标识为true
	if ep.CurrentPage == ep.Pages {
		ep.LastPage = true
	} else {
		ep.LastPage = false
	}

	return gr.Kvs, nil
}

//向前翻页
func (ep *EtcdPaginte) Previous() ([]*mvccpb.KeyValue, error) {

	if ep.TotalRecord == 0 {
		return nil, errors.New("no data returned")
	}

	if ep.CurrentPage == 1 {
		return nil, errors.New("no previous page")
	}

	opts := []clientv3.OpOption{
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(ep.PageSize),
		clientv3.WithFromKey(),
	}
	gr, err := ep.Kv.Get(ep.Ctx, ep.FirstKeyArray[ep.CurrentPage-2], opts...)
	if err != nil {
		return nil, err
	}
	ep.CurrentLastKey = string(gr.Kvs[len(gr.Kvs)-1].Key)
	ep.CurrentPage = ep.CurrentPage - 1

	//若已是最后一页则更改末页标识为true
	if ep.CurrentPage < ep.Pages {
		ep.LastPage = false
	}

	return gr.Kvs, nil

}
