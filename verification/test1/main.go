/**
 * @Author: DollarKillerX
 * @Description: main.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午3:49 2020/1/2
 */
package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

var (
	Etcd    *clientv3.Client // 客户端
	err     error
	Kv      clientv3.KV      // 用于读写etcd的kv
	Lease   clientv3.Lease   // 租约
	Watcher clientv3.Watcher // 监听
)

func main() {

	//test1()
	//leaseGrant()
	test3Watch()
}

func test1() {
	config := clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2079"},
		Username:    "golang",
		Password:    "123456",
		DialTimeout: 5 * time.Second,
	}

	Etcd, err = clientv3.New(config)
	if err != nil {
		panic(err.Error())
	}
	Kv = clientv3.NewKV(Etcd)
	Lease = clientv3.NewLease(Etcd)
	Watcher = clientv3.NewWatcher(Etcd)

	_, e := Kv.Put(context.TODO(), "/golang", "helloWorld")
	if e != nil {
		log.Println(e)
	}

	_, e = Kv.Put(context.TODO(), "/golang/ac", "helloWorld")
	if e != nil {
		log.Println(e)
	}

	kv, e := Kv.Get(context.TODO(), "/golang")
	if e != nil {
		log.Println(e)
	}
	for _, v := range kv.Kvs {
		log.Println(string(v.Key), "   :   ", string(v.Value))
	}

	kv, e = Kv.Get(context.TODO(), "/golang", clientv3.WithPrefix())
	if e != nil {
		log.Println(e)
	}
	for _, v := range kv.Kvs {
		log.Println(string(v.Key), "   :   ", string(v.Value))
	}
}

func leaseGrant() {
	config := clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2079"},
		Username:    "golang",
		Password:    "123456",
		DialTimeout: 5 * time.Second,
	}

	Etcd, err = clientv3.New(config)
	if err != nil {
		panic(err.Error())
	}
	Kv = clientv3.NewKV(Etcd)
	Lease = clientv3.NewLease(Etcd)
	Watcher = clientv3.NewWatcher(Etcd)

	les, e := Lease.Grant(context.TODO(), 10) // ttl 秒
	if e != nil {
		panic(e.Error())
	}

	// Put 一个KV 与租约关联实现10秒后过期
	leaseId := les.ID // 获取租约的id
	putResponse, e := Kv.Put(context.TODO(), "/cron/lock/job1", "SSr", clientv3.WithLease(leaseId))
	if e != nil {
		panic(e.Error())
	}

	// 自动续租
	//_, e = Lease.KeepAlive(context.TODO(), leaseId)
	//if e != nil {
	//	panic(e.Error())
	//}

	// 手动续租
	go func() {
		timer := time.NewTimer(time.Millisecond * 500)
		for {
			select {
			case <-timer.C:
				timeout, _ := context.WithTimeout(context.Background(), time.Second*3)
				_, e = Lease.KeepAlive(timeout, leaseId)
				if e != nil {
					panic(e.Error())
				}
			}
		}
	}()

	//timeout, _ := context.WithTimeout(context.Background(), time.Second*3)
	//_, e = Lease.KeepAlive(timeout, leaseId)
	//if e != nil {
	//	panic(e.Error())
	//}

	// 处理续租应答的协程
	//go func() {
	//forloop:
	//	for {
	//		select {
	//		case keepResp := <-etcdCh:
	//			if keepResp == nil {
	//				fmt.Println("租约已经失效")
	//				break forloop
	//			} else { // 每秒租约一次,所以就会受到一次应答
	//				fmt.Println("收到自动租约应答:", keepResp.ID)
	//			}
	//		}
	//	}
	//}()

	fmt.Println("写入成功", putResponse.Header.Revision)

	// 定时看一下kv过期没有
	for {
		getResponse, e := Kv.Get(context.TODO(), "/cron/lock", clientv3.WithPrefix())
		if e != nil {
			panic(e.Error())
		}

		if getResponse.Count == 0 {
			fmt.Println("Kv 过期")
			break
		}

		for _, v := range getResponse.Kvs {
			fmt.Println("还没有过期:", string(v.Key), " : ", string(v.Value))
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func test3Watch() {
	config := clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2079"},
		Username:    "golang",
		Password:    "123456",
		DialTimeout: 5 * time.Second,
	}

	client, e := clientv3.New(config)
	if e != nil {
		log.Fatalln(e)
	}

	kv := clientv3.NewKV(client)
	watcher := clientv3.NewWatcher(client)

	watch := watcher.Watch(context.TODO(), "/xxr", clientv3.WithPrefix())
	go func() {
		for {
			select {
			case <-watch:
				response, e := kv.Get(context.TODO(), "/xxr", clientv3.WithPrefix())
				if e != nil {
					log.Fatalln(e)
				}
				for _, v := range response.Kvs {
					log.Println(string(v.Key), "  :  ", string(v.Value))
				}
			}
		}
	}()

	for {
		time.Sleep(time.Second)
		kv.Put(context.TODO(), "/xxr/c", "sadsda")
	}
}
