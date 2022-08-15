package zetcd

import (
	"github.com/sirupsen/logrus"
)



type EtcdDis struct {
	Cluster string
	Etcd   *Etcd //etcd服务集群机器
	//注册; key:服务名称
	MapRegister map[string]*EtcdService  //["Cluster/server.name:Kind"]*EtcdService
	//监听相关的服务
	MapWatch map[string]*EtcdMaster //[service_name]*EtcdService
	Logger   *logrus.Logger
}

// NewEtcdDis 注册相关的函数
func NewEtcdDis(cluster string, etcd *Etcd, logger *logrus.Logger) *EtcdDis {
	return &EtcdDis{
		Cluster:     cluster,
		Etcd:       etcd,
		MapRegister: nil,
		MapWatch:    nil,
		Logger:      logger,
	}
}

func (d *EtcdDis) Register(keyInfo *KeyInfo, valueInfo *ValueInfo, retry bool) {
	var s *EtcdService
	var e error
	if s, e = NewEtcdService(d.Etcd, keyInfo, valueInfo, retry, d.Logger); e != nil {
		d.Logger.Warnf("[zetcd][Register] service:%s error:%s \n", s.Key.GetRegisterKey(), e.Error())
		return
	}
	if d.MapRegister == nil {
		d.MapRegister = make(map[string]*EtcdService)
	}
	if _, ok := d.MapRegister[s.Key.GetRegisterKey()]; ok {
		d.Logger.Warnf("[zetcd][Register] Service:%s Have Registered", s.Key.GetRegisterKey())
		return
	}
	d.MapRegister[s.Key.GetRegisterKey()] = s
	//维持心跳
	s.Start()
}

//监听相关的
func (d *EtcdDis) Watch(service string) {
	var w *EtcdMaster
	var e error
	if d.MapWatch == nil {
		d.MapWatch = make(map[string]*EtcdMaster)
	}
	if _, ok := d.MapWatch[service]; ok {
		d.Logger.Warnf("[zetcd][Watch] Service:%s Have Watch!\n", service)
		return
	}

	if w, e = NewMaster(d.Cluster,service,d.Etcd, d.Logger); e != nil {
		d.Logger.Warnf("[zetcd][Watch] Service:%s Failed! Error:%s\n", service, e.Error())
		return
	}
	go w.DoWatch()
	d.MapWatch[service] = w
}

func (d *EtcdDis) IsWatched(service string) bool {
	if _, ok := d.MapWatch[service]; ok {
		return true
	}
	return false
}

//GetGrpcServiceInfoRandom 获取服务grpc节点信息-随机获取
func (d *EtcdDis) GetGrpcServiceInfoRandom(service string) (*ValueInfo, bool) {
	if d.MapWatch == nil {
		d.Logger.Warnf("[zetcd][GetServiceInfoRandom] MapWatch is nil")
		return nil, false
	}
	if v, ok := d.MapWatch[service]; ok {
		if v != nil {
			if n, ok1 := v.GetGrpcNodeRandom(); ok1 {
				return n, true
			}
		}
	} else {
		d.Logger.Warnf("[zetcd][GetServiceInfoRandom] Service:%s Not Be Watched!\n", service)
	}
	return nil, false
}

//GetServiceInfoAllNode 获取服务的节点信息
func (d *EtcdDis) GetServiceInfoAllNode(service string) ([]EtcdNode, bool) {
	if d.MapWatch == nil {
		d.Logger.Warnf("[zetcd][GetServiceInfoAllNode] MapWatch is nil")
		return []EtcdNode{}, false
	}
	if v, ok := d.MapWatch[service]; ok {
		if v != nil {
			return v.GetAllNodes(), true
		}
	} else {
		d.Logger.Warnf("[zetcd][GetServiceInfoAllNode] Service:%s Not Be Watched!\n", service)
	}
	return []EtcdNode{}, false
}
