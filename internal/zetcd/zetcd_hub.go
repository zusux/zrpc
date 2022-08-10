package zetcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Hub struct {
	Hosts         []string
	DialTimeout   int64
	DialKeepAlive int64 //秒
	client        *clientv3.Client
}

func NewHub(hosts []string, dialTimeout,dialKeepAlive int64) (*Hub,error) {
	 h := Hub{
		Hosts: hosts,
		DialTimeout: dialTimeout,
		DialKeepAlive: dialKeepAlive,
	}
	err :=h.connect()
	return &h,err
}

// 设置etcd服务地址
func (z *Hub) SetEtcdServerAddress(etcdHosts []string) {
	z.Hosts = etcdHosts

}
// 添加etcd服务地址
func (z *Hub) AddEtcdServerAddress(addr []string) {
	z.Hosts = append(z.Hosts, addr...)
}

func (z *Hub) GetClient() *clientv3.Client {
	return z.client
}

func (z *Hub) connect() error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   z.Hosts,
		DialTimeout: time.Millisecond * time.Duration(z.DialTimeout),
	})
	z.client = client
	return err
}

func (z *Hub) put(key, value string, id clientv3.LeaseID) error {
	fmt.Printf("Register, Key:%s, Value:%s\n", key, value)
	_, err := z.client.Put(context.TODO(), key, value, clientv3.WithLease(id))
	return err
}


func (z *Hub) GetOne(key string, balance Balancer) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(z.DialTimeout))
	resp, err := z.client.Get(ctx, key)
	defer cancel()
	if err != nil {
		return "", err
	}
	idx, err := balance.GetPoint(len(resp.Kvs))
	if err != nil {
		return "", err
	}
	return string(resp.Kvs[idx].Value), nil
}

func (z *Hub) GetAll(key string) (*[]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(z.DialTimeout))
	resp, err := z.client.Get(ctx, key)
	defer cancel()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for idx := range resp.Kvs {
		result = append(result, string(resp.Kvs[idx].Value))
	}
	return &result, nil
}

func (z *Hub) grant() (*clientv3.LeaseGrantResponse, error) {
	return z.client.Grant(context.TODO(), z.DialKeepAlive)
}

func (z *Hub) revoke(id clientv3.LeaseID) error {
	_, err := z.client.Revoke(context.TODO(), id)
	return err
}

//LeaseTimeToLive 获取租约信息 todo
func (z *Hub) timeToLive(id clientv3.LeaseID) (*clientv3.LeaseTimeToLiveResponse, error) {
	return z.client.TimeToLive(context.TODO(), id)
}

func (z *Hub) keepAlive(id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return z.client.KeepAlive(context.TODO(), id)
}

func (z *Hub) watch(key string) {
	rch := z.client.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

//设置超时时间
func (z *Hub) SetDialTimeout(dialTimeout int64) *Hub {
	z.DialTimeout = dialTimeout
	return z
}
//设置keepalive时间
func (z *Hub) SetDialKeepAlive(dialKeepAlive int64) *Hub {
	z.DialKeepAlive = dialKeepAlive
	return z
}

