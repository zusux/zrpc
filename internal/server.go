package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/sd"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal/util"
	"github.com/zusux/zrpc/internal/zetcd"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

type Pub struct {
	Port    int
	Publish bool
}


type Server struct {
	Cluster string
	Ip     string
	Name    string
	Loger 	map[string]ZLog
	Etcd  zetcd.Etcd
	Retry  bool
	Listen map[string]Pub
	Mysql map[string]Mysql
	Redis map[string]Redis
	EtcdDis   *zetcd.EtcdDis
	Service  map[string]string //service.key, service.value
	sync.RWMutex
}

func NewServer() (*Server, error) {
	siteMode := GetSiteMode()
	LoadEnv(siteMode)
	svr := Server{}
	// Quick unmarshal.
	err := K.Unmarshal("server", &svr)
	if err != nil {
		log.Printf("[zrpc][error] unmarshal config err: %s",err)
		return &svr, err
	}
	svr.FillIp()
	svr.SetSiteMode(siteMode)
	svr.Log()
	return &svr, nil
}

func (s *Server) FillIp()  {
	if s.Ip == "" {
		ip, err := util.GetLocalIP()
		if err != nil {
			log.Printf("[zrpc] reg get local ip err: %s",err)
		}
		s.Ip = ip.String()
	}
}
func (s *Server) SetSiteMode(siteMode string)  {
	s.Cluster = siteMode
}

func GetSiteMode() string {
	siteMode := os.Getenv("site_mode")
	if siteMode == ""{
		siteMode = "prod"
	}
	return siteMode
}

func (s *Server) Log() *logrus.Logger {
	return s.GetLog("default_log")
}

func (s *Server) GetLog(section string) *logrus.Logger {
	s.RLock()
	defer s.RUnlock()
	section = s.getSectionLogName(section)
	zLog,ok := s.Loger[section]
	if ok{
		if !zLog.Initialize{
			zLog.init()
		}
		return zLog.Log
	}
	return logrus.New()
}

func (s *Server) getSectionLogName(section string) string{
	for k,_ := range s.Loger{
		if section == k{
			return section
		}
	}
	return "default_log"
}

func (s *Server) Register() {
	//注册服务
	if s.EtcdDis == nil {
		s.EtcdDis = zetcd.NewEtcdDis(s.Cluster, &s.Etcd, s.Log())
	}
	for kind, v := range s.Listen {
		if v.Publish {
			keyInfo := zetcd.NewKeyInfo(s.Cluster, s.Name, kind)
			keyValue := zetcd.NewValueInfo(kind, s.Ip, uint32(v.Port), zetcd.StatusInit, 0, time.Now().Unix())
			s.EtcdDis.Register(keyInfo, keyValue, s.Retry)
		}
	}
	//监测service
	s.EtcdDis.Watch(s.Name)
	for _, Svr := range s.Service {
		s.EtcdDis.Watch(Svr)
	}
}
func (s *Server) UnRegister() {
	//关闭检查的依赖服务
	for key, service := range s.EtcdDis.MapWatch {
		service.Close()
		s.Log().Warnf("[zrpc] Service %s unRegister", key)
	}
	//关闭注册服务
	for key, server := range s.EtcdDis.MapRegister {
		server.Stop()
		s.Log().Warnf("[zrpc] Server %s unRegister", key)
	}


}

// GrpcRequestRemote 外部服务
func (s *Server) GrpcRequestRemote(ctx context.Context, serverName string, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	//创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	if !s.EtcdDis.IsWatched(serverName) {
		s.EtcdDis.Watch(serverName)
	}
	valueInfo, ok := s.EtcdDis.GetGrpcServiceInfoRandom(serverName)
	if ok {
		addr, ok := valueInfo.GetValidAddress()
		if ok {
			ed, _, _ := reqFactory(addr)
			response, err = ed(ctx, req)
		} else {
			s.Log().Errorf("[zetcd][service] getValidAddress addr:%s, ok:%v", addr, ok)
			err = errors.New(serverName+ " service endpoint address is not valid")
		}
	} else {
		s.Log().Errorf("[zetcd][service] %s have not grpc endpoint:%v, ok:%v",serverName, valueInfo, ok)
		err = errors.New(serverName+ " service no find endpoint")
	}
	return
}

// GetRandomEndPoint 获取随机的端点
func (s *Server) GetRandomEndPoint(serviceName string) (endpoint string, err error) {
	valueInfo, ok := s.EtcdDis.GetGrpcServiceInfoRandom(serviceName)
	if ok {
		addr, ok := valueInfo.GetValidAddress()
		if ok {
			return addr,nil
		} else {
			s.Log().Errorf("[zetcd][service] service endpoint address is not valid addr:%s, ok:%v", addr, ok)
			err = errors.New(serviceName+ " service endpoint address is not valid")
		}
	} else {
		s.Log().Errorf("[zetcd][service] %s have not endpoint",serviceName)
		err = errors.New(serviceName+ " service no find endpoint")
	}
	return
}

// GetDb 指定DB
func (s *Server)GetDb() (*gorm.DB, error) {
	return s.GetDbBySec("default_mysql")
}

// GetDbBySec 指定DB
func (s *Server)GetDbBySec(sectionDb string) (*gorm.DB, error) {
	mysql,ok := s.Mysql[sectionDb]
	if ok {
		log := s.GetLog(mysql.SectionLog)
		mysql.SetLogger(log)
		return mysql.GetDb()
	}
	return nil, errors.New(fmt.Sprintf("mysql config %s is not find",sectionDb))
}