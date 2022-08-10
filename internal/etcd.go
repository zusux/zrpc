package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/sd"
	"github.com/zusux/zrpc/internal/zetcd"
	"time"
)

type Pub struct {
	Kind    string
	Port    int  `koanf:"port"`
	Publish bool `koanf:"publish"`
}

type Reg struct {
	Cluster       string
	Name          string
	Ip            string
	Hosts         []string
	DialTimeout   int64
	DialKeepAlive int64
	Retry         bool
	Pubs          map[string]Pub
	EtcdDis       *zetcd.EtcdDis
	server          *Server
}

func NewReg(cluster, name, ip string, dialTimeout, dialKeepAlive int64, retry bool, host []string, pubs map[string]Pub) *Reg {
	return &Reg{
		Cluster:       cluster,
		Name:          name,
		Ip:            ip,
		Hosts:         host,
		DialTimeout:   dialTimeout,
		DialKeepAlive: dialKeepAlive,
		Retry:         retry,
		Pubs:          pubs,
	}
}

func (s *Reg) Register() {
	if s.EtcdDis == nil {
		s.EtcdDis = zetcd.NewEtcdDis(s.Cluster, s.Hosts, Log())
	}
	for _, v := range s.Pubs {
		if v.Publish {
			keyinfo := zetcd.NewKeyInfo(s.Cluster, s.Name, v.Kind)
			keyValue := zetcd.NewValueInfo(v.Kind, s.Ip, uint32(v.Port), zetcd.StatusInit, 0, time.Now().Unix())
			s.EtcdDis.Register(keyinfo, keyValue, s.DialTimeout, s.DialKeepAlive, s.Retry)
		}
	}
}
func (s *Reg) UnRegister() {
	for key, server := range s.EtcdDis.MapRegister {
		err := server.Revoke()
		if err != nil {
			Log().Error(fmt.Sprintf("[zetcd][core] %s UnRegister Error, %s", key, err))
		}
	}
}

// GrpcRequestRemote 外部服务
func (s *Reg) GrpcRequestRemote(ctx context.Context, serverName string, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	if !s.EtcdDis.IsWatched(serverName) {
		s.EtcdDis.Watch(serverName)
	}
	etcdNode, ok := s.EtcdDis.GetServiceInfoRandom(serverName)
	if ok {
		addr, ok := etcdNode.GetValidAddress()
		if ok {
			ed, _, _ := reqFactory(addr)
			response, err = ed(ctx, req)
		} else {
			Log().Error(fmt.Sprintf("[zetcd][core] getValidAddress addr:%s, ok:%v", addr, ok))
		}
	} else {
		Log().Error(fmt.Sprintf("[zetcd][core] getValidAddress etcdNode:%v, ok:%v", etcdNode, ok))
		err = errors.New("no find node")
	}
	return
}
