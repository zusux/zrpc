package zetcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"math/rand"
	"sync"
)

type EtcdMaster struct {
	Cluster string
	Etcd    *Etcd
	Service string               //路径
	Nodes   map[string]*EtcdNode //[grpc]*EtcdNode
	hub     *Hub
	stop    chan error
	State   bool
	Loger   *logrus.Logger
	sync.Mutex
}

//EtcdNode etcd注册的节点
type EtcdNode struct {
	State bool
	*KeyInfo
	EndPoints map[string]*ValueInfo //[ip:port]ValueInfo
	sync.Mutex
}

func (m *EtcdMaster) NewEtcdNodes(state bool, keyInfo *KeyInfo) *EtcdNode {
	return &EtcdNode{
		State:     state,
		KeyInfo:   keyInfo,
		EndPoints: make(map[string]*ValueInfo),
	}
}

//添加节点
func (n *EtcdNode) addEndPoint(valueInfo *ValueInfo) {
	n.Lock()
	defer n.Unlock()
	n.EndPoints[valueInfo.getRegisterAddress()] = valueInfo
}

//删除节点
func (n *EtcdNode) deleteEndPoint(valueInfo *ValueInfo) {
	n.Lock()
	defer n.Unlock()
	_, ok := n.EndPoints[valueInfo.getRegisterAddress()]
	if ok {
		delete(n.EndPoints, valueInfo.getRegisterAddress())
	}
}

func (n *EtcdNode) GetEndPointRandom() (*ValueInfo, bool) {
	n.Lock()
	defer n.Unlock()
	count := len(n.EndPoints)
	if count == 0 {
		return nil, false
	}
	idx := rand.Intn(count)
	for _, v := range n.EndPoints {
		if idx == 0 {
			return v, true
		}
		idx = idx - 1
	}
	return nil, false
}

func NewMaster(cluster, service string, etcd *Etcd, logger *logrus.Logger) (*EtcdMaster, error) {
	hub, err := NewHub(etcd)
	master := &EtcdMaster{
		Cluster: cluster,
		Etcd:    etcd,
		Service: service,
		Nodes:   make(map[string]*EtcdNode),
		hub:     hub,
		Loger:   logger,
		stop:    make(chan error),
		State:   true,
	}
	return master, err
}

func (m *EtcdMaster) GetPrefixPath() string {
	return fmt.Sprintf("%s/%s", m.Cluster, m.Service)
}

func (m *EtcdMaster) DecodeKeyValue(ev *mvccpb.KeyValue) (*KeyInfo, *ValueInfo, error) {
	m.Loger.Infof("[zetcd][DecodeKeyValue] decode key:%q,value:%q\n", ev.Key, ev.Value)
	keyInfo, err1 := m.GetKeyInfo(ev.Key)
	if err1 != nil {
		m.Loger.Errorf("[zetcd][DecodeKeyValue] %v decode key error: %s", string(ev.Key), err1)
		return nil, nil, err1
	}
	valueInfo, err2 := m.GetValueInfo(ev.Value)
	if err2 != nil {
		m.Loger.Errorf("[zetcd][DecodeKeyValue] %v decode value error: %s", string(ev.Value), err2)
		return nil, nil, err2
	}
	return keyInfo, valueInfo, nil
}

func (m *EtcdMaster) DoWatch() {
	m.Get()
	m.Watch()
}

func (m *EtcdMaster) Get() {
	//查看之前存在的节点
	resp, err := m.hub.GetClient().Get(context.Background(), m.GetPrefixPath(), clientv3.WithPrefix())
	if err != nil {
		m.Loger.Errorf("[zetcd][WatchNodes] get Node error! prefixPath:%s error:%s\n", m.GetPrefixPath(), err.Error())
	} else {
		for _, ev := range resp.Kvs {
			keyInfo, valueInfo, err := m.DecodeKeyValue(ev)
			if err != nil {
				continue
			}
			m.addNode(keyInfo, valueInfo)
		}
	}
}

func (m *EtcdMaster) ReBuildNode(keyInfo *KeyInfo) {
	//查看之前存在的节点
	resp, err := m.hub.GetClient().Get(context.Background(), keyInfo.GetRegisterKey(), clientv3.WithPrefix())
	if err != nil {
		m.Loger.Errorf("[zetcd][WatchNodes] get Node error! prefixPath:%s error:%s\n", m.GetPrefixPath(), err.Error())
		return
	}
	newMapNode := make(map[string]*EtcdNode, 0)
	for _, ev := range resp.Kvs {
		keyInfo, valueInfo, err := m.DecodeKeyValue(ev)
		if err != nil {
			continue
		}
		node := m.NewEtcdNodes(true, keyInfo)
		node.addEndPoint(valueInfo)
		newMapNode[keyInfo.Kind] = node
	}
	m.Lock()
	defer m.Unlock()
	m.Nodes[keyInfo.Kind] = newMapNode[keyInfo.Kind]
}

func (m *EtcdMaster) Watch() {

	rch := m.hub.GetClient().Watch(context.Background(), m.GetPrefixPath(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for {
		select {
		case <-m.stop:
			return
		case wResp := <-rch:
			for _, ev := range wResp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					m.Loger.Infof("[zetcd][watch] type:%s  key:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					keyInfo, valueInfo, err := m.DecodeKeyValue(ev.Kv)
					if err != nil {
						m.Loger.Errorf("[zetcd][watch] type:[%s] key:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
						continue
					}
					m.addNode(keyInfo, valueInfo)
				case clientv3.EventTypeDelete:
					m.Loger.Warnf("[zetcd][watch] type:[%s] key:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					if m.State {
						keyInfo, _, err := m.DecodeKeyValue(ev.Kv)
						if err != nil {
							m.Loger.Errorf("[zetcd][watch] type:[%s] key:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
							continue
						}
						m.ReBuildNode(keyInfo)
					}
				default:
					m.Loger.Infof("[zetcd][watch] type:[%s] key:%q, value:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				}
			}
		}
	}
}

func (m *EtcdMaster) GetKeyInfo(key []byte) (*KeyInfo, error) {
	keyInfo := &KeyInfo{}
	err := keyInfo.SetRegisterKey(string(key))
	return keyInfo, err
}

func (m *EtcdMaster) GetValueInfo(value []byte) (*ValueInfo, error) {
	valueInfo := &ValueInfo{}
	if len(value) == 0 {
		return valueInfo, errors.New(fmt.Sprintf("got empty value %s", value))
	}
	err := valueInfo.DecodeValue(value)
	return valueInfo, err
}

//添加节点
func (m *EtcdMaster) addNode(keyInfo *KeyInfo, valueInfo *ValueInfo) {
	node, ok := m.Nodes[keyInfo.Kind]
	if !ok {
		m.Lock()
		defer m.Unlock()
		node = m.NewEtcdNodes(true, keyInfo)
		m.Nodes[keyInfo.Kind] = node
	}
	node.addEndPoint(valueInfo)
}

//删除节点
func (m *EtcdMaster) deleteNode(keyInfo *KeyInfo, valueInfo *ValueInfo) {
	node, ok := m.Nodes[keyInfo.Kind]
	if ok {
		node.deleteEndPoint(valueInfo)
	}
	if len(node.EndPoints) == 0 {
		m.Lock()
		defer m.Unlock()
		delete(m.Nodes, keyInfo.Kind)
	}
}

//deleteEndPoint 删除端点
func (m *EtcdMaster) deleteEndPoint(key string, keyInfo *KeyInfo, valueInfo *ValueInfo) {
	node, ok := m.Nodes[keyInfo.Kind]
	if !ok {
		node = m.NewEtcdNodes(true, keyInfo)
		m.Nodes[key] = node
	}
	node.addEndPoint(valueInfo)
}

func (m *EtcdMaster) GetAllNodes() []EtcdNode {
	var temp []EtcdNode
	for _, v := range m.Nodes {
		if v != nil {
			temp = append(temp, *v)
		}
	}
	return temp
}

func (m *EtcdMaster) GetGrpcNodeRandom() (*ValueInfo, bool) {
	etcdNode, ok := m.Nodes["grpc"]
	if ok {
		return etcdNode.GetEndPointRandom()
	}
	return nil, false
}
func (m *EtcdMaster) Stop() {
	m.stop <- nil
}
func (m *EtcdMaster) Close() {
	m.Stop()
}
