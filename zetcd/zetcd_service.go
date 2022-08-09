package zetcd

import (
	"fmt"
	"github.com/sirupsen/logrus"
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

func NewEtcdService(hosts []string, keyinfo *KeyInfo, valueInfo *ValueInfo, dialTimeout, dialKeepAlive int64, retry bool, log *logrus.Logger) (*EtcdService, error) {

	hub,err := NewHub(hosts,dialTimeout,dialKeepAlive)
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
		s.logger.WithField("init", "start").Error(fmt.Sprintf("[zetcd] client keeplive fail err:%s , code: %d \n", err, s.status))
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
				s.logger.WithField("keepalive", "client.Done").Info(fmt.Sprintf("[zetcd] server closed  \n"))
				s.analyseRetry()
				return
			case res, ok := <-ch:
				if !ok {
					s.logger.WithField("keepalive", "close").Info(fmt.Sprintf("[zetcd] keep live channel info res:%v ok:%v \n", res, ok))
					s.logger.WithField("keepalive", "close").Info(fmt.Sprintf("[zetcd] keep live channel closed code:%d \n", s.status))
					s.Revoke()
					s.analyseRetry()
					return
				} else {
					s.logger.WithField("keepalive", "online").Info(fmt.Sprintf("[zetcd] recv reply from service:%s, ttl:%d\n", s.Key.GetRegisterKey(), res.TTL))
				}
			}
		}
	}()
	return nil
}


func (s *EtcdService) analyseRetry() {
	if s.Retry && s.status >= 500 {
		s.logger.WithField("keepalive", "analyseRetry").Info(fmt.Sprintf("[zetcd] start to restart code: %d \n", s.status))
		s.ReStart()
	}else{
		s.logger.WithField("keepalive", "analyseRetry").Info(fmt.Sprintf("[zetcd]  restart do nothing code: %d \n", s.status))
	}
}



func (s *EtcdService) ReStart() error {
	err := s.hub.connect()
	if err != nil {
		s.status = StatusConnectError
		s.logger.Error(fmt.Sprintf("[zetcd] client reconnect fail err:%s , code: %d \n", err, s.status))
		return err
	}
	err = s.Revoke()
	if err != nil {
		s.status = StatusRevokeError
		s.logger.Error(fmt.Sprintf("[zetcd] client Revoke fail err:%s , code: %d \n", err, s.status))
	}
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
		s.logger.Error(fmt.Sprintf("[zetcd] client update service fail err:%s , value: %v \n", err, s.Value))
		return err
	}
	return nil
}

//KeepLive 用于维持租约
func (s *EtcdService) KeepLive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	resp, err := s.hub.grant()
	if err != nil {
		s.status = StatusGrantError
		s.logger.Error(fmt.Sprintf("[zetcd] client grant fail err:%s , code: %d \n", err, s.status))
		return nil, err
	}
	valueByte,err := s.Value.EncodeValue()
	err = s.hub.put(s.Key.GetRegisterKey(),string(valueByte) , resp.ID)
	if err != nil {
		s.status = StatusPutError
		s.logger.Error(fmt.Sprintf("[zetcd] client put fail err:%s , code: %d \n", err, s.status))
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
func (s *EtcdService) Revoke() error {
	err := s.hub.revoke(s.leaseId)
	if err != nil {
		s.status = StatusRevokeError
		s.logger.Error(fmt.Sprintf("[zetcd] client Revoke fail err:%s , code: %d\n", err, s.status))
	}
	s.logger.Info(fmt.Sprintf("[zetcd] service: %s stop\n", s.Key.GetRegisterKey()))
	return err
}
