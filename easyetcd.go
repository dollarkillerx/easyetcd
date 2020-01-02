/**
 * @Author: DollarKillerX
 * @Description: easyetcd.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午8:34 2020/1/2
 */
package easyetcd

import "github.com/coreos/etcd/clientv3"

type easyEtcd struct {
	config clientv3.Config
	client *clientv3.Client
}

func New(config clientv3.Config) (*easyEtcd, error) {
	client, e := clientv3.New(config)
	if e != nil {
		return nil, e
	}
	return &easyEtcd{
		config: config,
		client: client,
	}, nil
}

func (e *easyEtcd) NewKv() clientv3.KV {
	return clientv3.NewKV(e.client)
}

func (e *easyEtcd) NewLease() clientv3.Lease {
	return clientv3.NewLease(e.client)
}

func (e *easyEtcd) NewWatcher() clientv3.Watcher {
	return clientv3.NewWatcher(e.client)
}

