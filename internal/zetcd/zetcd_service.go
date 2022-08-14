package zetcd

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

//STATUS 0 初始化, 1 正常服务, 2,挂起, 3, 主动关闭 500 异常情况
type STATUS int

const StatusInit STATUS = 0
const StatusOk STATUS = 1
const StatusHang STATUS = 2
const StatusClose STATUS = 3
const StatusConnectError STATUS = 500
const StatusGrantError STATUS = 501
const StatusPutError STATUS = 502
const StatusKeepaliveError STATUS = 503
const StatusRevokeError STATUS = 504
const StatusEtcdRemoteCloseError STATUS = 505

type EtcdService struct {
	Key  *KeyInfo
	Value  *ValueInfo
	Retry   bool
	stop    chan error
	leaseId clientv3.LeaseID
	hub  *Hub
	status  STATUS //0 初始化, 1 正常服务, 2 挂起, 3 主动关闭 500 异常情况,
	logger  *logrus.Logger
}

func NewEtcdService(etcd *Etcd, keyinfo *KeyInfo, valueInfo *ValueInfo, retry bool, log *logrus.Logger) (*EtcdService, error) {

	hub,err := NewHub(etcd)
	service := EtcdService{
		Key:   keyinfo,
		Value:   valueInfo,
		Retry:  retry,
		stop:   make(chan error),
		hub: hub,
		status: StatusInit,
		logger: log,
	}
	if err == nil {
		service.status = StatusOk
	}
	return &service, err
}

// Start 启动
func (s *EtcdService) Start() error {
	ch, err := s.KeepLive()
	if err != nil {
		s.status = StatusKeepaliveError
		s.logger.WithField("init", "start").Errorf("[zetcd] client keeplive fail err:%s , code: %d \n", err, s.status)
		return err
	}
	go func() {
		for {
			select {
			case <-s.stop:
				//主动停止
				s.Revoke()
				return
			case <-s.hub.client.Ctx().Done():
				s.logger.WithField("keepalive", "client.Done").Infof("[zetcd] server closed  \n")
				s.analyseRetry()
				return
			case res, ok := <-ch:
				if !ok {
					s.logger.WithField("keepalive", "close").Infof("[zetcd] keep live channel info res:%v ok:%v \n", res, ok)
					s.logger.WithField("keepalive", "close").Infof("[zetcd] keep live channel closed code:%d \n", s.status)
					s.status = StatusEtcdRemoteCloseError
					s.Revoke()
					s.analyseRetry()
					return
				} else {
					s.logger.WithField("keepalive", "online").Infof("[zetcd] recv reply from service:%s, ttl:%d\n", s.Key.GetRegisterKey(), res.TTL)
				}
			}
		}
	}()
	return nil
}


func (s *EtcdService) analyseRetry() {
	if s.Retry && s.status >= 500 {
		s.logger.WithField("keepalive", "analyseRetry").Infof("[zetcd] start to restart code: %d \n", s.status)
		s.ReStart()
	}else{
		s.logger.WithField("keepalive", "analyseRetry").Infof("[zetcd]  restart do nothing code: %d \n", s.status)
	}
}



func (s *EtcdService) ReStart() error {
	err := s.hub.connect()
	if err != nil {
		s.status = StatusConnectError
		s.logger.Errorf("[zetcd] client reconnect fail err:%s , code: %d \n", err, s.status)
		return err
	}
	s.Revoke()
	s.Start()
	return nil
}

func (s *EtcdService) UpdateService(status STATUS, requestFlow uint32) error{
	s.Value.Status = status
	s.Value.RequestFlow = requestFlow

	valueByte,err := s.Value.EncodeValue()
	err = s.hub.put(s.Key.GetRegisterKey(),string(valueByte) , s.leaseId)
	if err != nil {
		s.status = StatusPutError
		s.logger.Errorf("[zetcd] client update service fail err:%s , value: %v \n", err, s.Value)
		return err
	}
	return nil
}
func (s *EtcdService) GetInfo() ([]*mvccpb.KeyValue,error){
	//查看之前存在的节点
	resp, err := s.hub.GetClient().Get(context.Background(), s.Key.GetRegisterKey(), clientv3.WithPrefix())
	if err != nil {
		s.logger.Errorf("[zetcd][WatchNodes] get Node error! prefixPath:%s error:%s\n", s.Key.GetRegisterKey(), err.Error())
	}
	return resp.Kvs,err
}


//KeepLive 用于维持租约
func (s *EtcdService) KeepLive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	resp, err := s.hub.grant()
	if err != nil {
		s.status = StatusGrantError
		s.logger.Errorf("[zetcd] client grant fail err:%s , code: %d \n", err, s.status)
		return nil, err
	}
	valueByte,err := s.Value.EncodeValue()
	err = s.hub.put(s.Key.GetRegisterKey(),string(valueByte) , resp.ID)
	if err != nil {
		s.status = StatusPutError
		s.logger.Errorf("[zetcd] client put %s fail err:%s , code: %d \n",s.Key.GetRegisterKey(), err, s.status)
		return nil, err
	}
	s.leaseId = resp.ID
	return s.hub.keepAlive(resp.ID)
}


// Stop 停止
func (s *EtcdService) Stop() {
	s.status = StatusClose
	s.stop <- nil
}

//revoke 撤销一个租约
func (s *EtcdService) Revoke(){
	go func() {
		err := s.hub.revoke(s.leaseId)
		if err != nil {
			s.status = StatusRevokeError
			s.logger.Errorf("[zetcd] client Revoke fail err:%s , code: %d\n", err, s.status)
		}
		s.logger.Infof("[zetcd] service: %s stop\n", s.Key.GetRegisterKey())
	}()
}