package zetcd

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type EtcdDis struct {
	Cluster string
	Hosts   []string //etcd服务集群机器
	//注册; key:服务名称
	MapRegister map[string]*EtcdService
	//监听相关的服务
	MapWatch map[string]*EtcdMaster
	logger   *logrus.Logger
}

// NewEtcdDis 注册相关的函数
func NewEtcdDis(cluster string, hosts []string, logger *logrus.Logger) *EtcdDis {
	return &EtcdDis{
		Cluster:     cluster,
		Hosts:       hosts,
		MapRegister: nil,
		MapWatch:    nil,
		logger:      logger,
	}
}

func (d *EtcdDis) Register(keyInfo *KeyInfo, valueInfo *ValueInfo, dialTimeout, dialKeepAlive int64, retry bool) {

	var s *EtcdService
	var e error
	if s, e = NewEtcdService(d.Hosts, keyInfo, valueInfo, dialTimeout, dialKeepAlive, retry, d.logger); e != nil {
		fmt.Printf("[zetcd] Register service:%s error:%s \n", s.Key.GetRegisterKey(), e.Error())
		return
	}
	if d.MapRegister == nil {
		d.MapRegister = make(map[string]*EtcdService)
	}
	if _, ok := d.MapRegister[s.Key.GetRegisterKey()]; ok {
		fmt.Printf("Service:%s Have Registered", s.Key.GetRegisterKey())
		return
	}
	d.MapRegister[s.Key.GetRegisterKey()] = s
	//维持心跳
	s.Start()
	d.Watch(s.Key.GetRegisterKey())
}

//监听相关的
func (d *EtcdDis) Watch(service string) {
	var w *EtcdMaster
	var e error
	if w, e = NewMaster(d.Hosts, service); e != nil {
		fmt.Printf("Watch Service:%s Failed! Error:%s\n", service, e.Error())
		return
	}
	if d.MapWatch == nil {
		d.MapWatch = make(map[string]*EtcdMaster)
	}

	if _, ok := d.MapWatch[service]; ok {
		fmt.Printf("Service:%s Have Watch!\n", service)
		return
	}
	d.MapWatch[service] = w
}

func (d *EtcdDis) IsWatched(service string) bool {
	if _, ok := d.MapWatch[service]; ok {
		return true
	}
	return false
}

//GetServiceInfoRandom 获取服务节点信息-随机获取
func (d *EtcdDis) GetServiceInfoRandom(service string) (EtcdNode, bool) {
	if d.MapWatch == nil {
		fmt.Println("MapWatch is nil")
		return EtcdNode{}, false
	}
	if v, ok := d.MapWatch[service]; ok {
		if v != nil {
			if n, ok1 := v.GetNodeRandom(); ok1 {
				return n, true
			}
		}
	} else {
		fmt.Printf("Service:%s Not Be Watched!\n", service)
	}
	return EtcdNode{}, false
}

//GetServiceInfoAllNode 获取服务的节点信息
func (d *EtcdDis) GetServiceInfoAllNode(service string) ([]EtcdNode, bool) {
	if d.MapWatch == nil {
		fmt.Println("MapWatch is nil")
		return []EtcdNode{}, false
	}
	if v, ok := d.MapWatch[service]; ok {
		if v != nil {
			return v.GetAllNodes(), true
		}
	} else {
		fmt.Printf("Service:%s Not Be Watched!\n", service)
	}
	return []EtcdNode{}, false
}
