package zetcd

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"github.com/zusux/zrpc/env"
	"os"
	"strings"
	"time"
)

type Etcd struct {
	ServerName        string // "/services/book/"
	Kind    string  //grpc,http
	Port int
	ServerAddress string `mapstructure:"server_address"` //"127.0.0.1:5040"
	EtcdServerAddress []string      `mapstructure:"etcd_server_address"` //"127.0.0.1:2379"
	DialTimeout       time.Duration `mapstructure:"dial_timeout"`
	DialKeepAlive     time.Duration `mapstructure:"dial_keep_alive"`

	service  etcdv3.Service
	client etcdv3.Client
}

func NewEtcd(serverName, kind,serverAddress string,  dialTimeout, dialKeepAlive int64 ,etcdAddr ...string) *Etcd {
	return &Etcd{
		ServerName:            serverName,
		Kind: kind,
		ServerAddress:     serverAddress,
		EtcdServerAddress: etcdAddr,
		DialTimeout:       time.Duration(dialTimeout) ,
		DialKeepAlive:     time.Duration(dialKeepAlive),
	}
}

func NewEtcdForClient() *Etcd {
	return NewEtcd(
		strings.TrimRight(env.K.String("server.name"), "/"),
		"",
		"",
		env.K.Int64("etcd.dial_timeout"),
		env.K.Int64("etcd.dial_keep_alive"),
		env.K.Strings("etcd.etcd_server_address")...
	)
}

func (e *Etcd) GetClient() etcdv3.Client {
	if e.client != nil {
		return e.client
	}
	var ctx = context.Background()
	//etcd连接参数
	option := etcdv3.ClientOptions{
		DialTimeout: time.Millisecond * time.Duration(e.DialTimeout),
		DialKeepAlive: time.Millisecond * time.Duration(e.DialKeepAlive),
	}
	//创建连接
	client, err := etcdv3.NewClient(ctx, e.EtcdServerAddress, option)
	if err != nil {
		panic(err)
	}
	e.client = client
	return client
}
// 设置etcd服务地址
func (e *Etcd) SetEtcdServerAddress(etcdServerAddresss []string) *Etcd {
	e.EtcdServerAddress = etcdServerAddresss
	return e
}
// 添加etcd服务地址
func (e *Etcd) AddEtcdServerAddress(addr ...string) *Etcd {
	e.EtcdServerAddress = append(e.EtcdServerAddress, addr...)
	return e
}

//设置路径
func (e *Etcd) SetPrefix(serverName string) *Etcd {
	e.ServerName = serverName
	return e
}


//设置超时时间
func (e *Etcd) SetDialTimeout(dialTimeout int) *Etcd {
	e.DialTimeout = time.Duration(dialTimeout)
	return e
}

//设置keepalive时间
func (e *Etcd) SetDialKeepAlive(dialKeepAlive int) *Etcd {
	e.DialKeepAlive = time.Duration(dialKeepAlive)
	return e
}

//设置服务地址 ip:port
func (e *Etcd) SetServerAddress(serverAddresss string) *Etcd {
	e.ServerAddress = serverAddresss
	return e
}
//注册
func (e *Etcd) Register() {
	e.service = etcdv3.Service{
		Key: e.GetRegPrefix(e.ServerName,e.Kind),
		Value: e.ServerAddress,
	}

	//创建注册
	registrar := etcdv3.NewRegistrar(e.GetClient(), e.service, log.NewJSONLogger(os.Stdout))
	registrar.Register() //启动注册服务
}

//获取注册的key
func (e *Etcd) GetRegPrefix(serverName,kind string) string {
	// key "/service/foobar/grpc:8080"
	return fmt.Sprintf("%s:%s",serverName,kind)
}

//反注册
func (e *Etcd) UnRegister() error {
	return e.client.Deregister(e.service)
}

//外部服务 grpc请求
func (e *Etcd) RequestRemote(ctx context.Context, serverName string,kind string, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	logger := log.NewNopLogger()
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	instancer, err := etcdv3.NewInstancer(e.GetClient(), e.GetRegPrefix(serverName,kind), logger)
	if err != nil {
		return nil,err
	}
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger) //reqFactory自定义的函数，主要用于端点层（endpoint）接受并显示数据
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	我们可以通过负载均衡器直接获取请求的endPoint，发起请求
	reqEndPoint,_ := balancer.Endpoint()
	*/

	/**
	也可以通过retry定义尝试次数进行请求
	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

	//现在我们可以通过 endPoint 发起请求了
	response, err = reqEndPoint(ctx, req)
	return
}


// 内部服务 grpc请求
func (e *Etcd) RequestLocal(ctx context.Context, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	logger := log.NewNopLogger()
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	instancer, err := etcdv3.NewInstancer(e.GetClient(), e.GetRegPrefix(e.ServerName,e.Kind), logger)
	if err != nil {
		return nil,err
	}
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger) //reqFactory自定义的函数，主要用于端点层（endpoint）接受并显示数据
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	/**
	我们可以通过负载均衡器直接获取请求的endPoint，发起请求
	reqEndPoint,_ := balancer.Endpoint()
	*/

	/**
	也可以通过retry定义尝试次数进行请求
	*/
	reqEndPoint := lb.Retry(3, 3*time.Second, balancer)

	//现在我们可以通过 endPoint 发起请求了
	response, err = reqEndPoint(ctx, req)
	return
}
