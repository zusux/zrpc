package zetcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"math/rand"
)

type EtcdMaster struct {
	Path string //路径
	Nodes map[string]*EtcdNode
	hub *Hub
}

//EtcdNode etcd注册的节点
type EtcdNode struct{
	State bool
	*KeyInfo
	*ValueInfo
}

func NewMaster(hosts []string,watchPath string)(*EtcdMaster,error){
	hub,err := NewHub(hosts,1,1)
	master := &EtcdMaster{
		Path:    watchPath,
		Nodes:   make(map[string]*EtcdNode),
		hub:  hub,
	}
	go master.WatchNodes()
	return master,err
}

func NewEtcdNode(state bool, keyInfo *KeyInfo, valueInfo *ValueInfo) *EtcdNode  {
	return &EtcdNode{
		State: state,
		KeyInfo:keyInfo,
		ValueInfo:valueInfo,
	}
}


func (m *EtcdMaster) WatchNodes(){
	//查看之前存在的节点
	resp,err := m.hub.GetClient().Get(context.Background(),m.Path,clientv3.WithPrefix())
	if err != nil{
		fmt.Printf("[zetcd] get nodes error:%s\n",err.Error())
	}else{
		for _,ev := range resp.Kvs{
			fmt.Printf("[zetcd] add dir:%q,value:%q\n",ev.Key,ev.Value)
			keyInfo,err1 := m.GetKeyInfo(ev.Key)
			if err1 != nil{
				fmt.Println(err1)
				continue
			}
			valueInfo,err2 := m.GetValueInfo(ev.Value)
			if err2 != nil{
				fmt.Println(err2)
				continue
			}
			m.addNode(string(ev.Value),keyInfo,valueInfo)
		}
	}

	rch := m.hub.GetClient().Watch(context.Background(),m.Path,clientv3.WithPrefix(),clientv3.WithPrevKV())
	for wResp := range rch{
		for _,ev := range wResp.Events{
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] dir:%q, value:%q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
				keyInfo,err1 := m.GetKeyInfo(ev.Kv.Key)
				valueInfo,err2 := m.GetValueInfo(ev.Kv.Value)
				if err1 != nil {
					fmt.Println(err1)
					continue
				}
				if err2 != nil {
					fmt.Println(err2)
					continue
				}
				m.addNode(string(ev.Kv.Key),keyInfo,valueInfo)
			case clientv3.EventTypeDelete:
				fmt.Printf("[%s] dir:%q, value:%q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
				k := ev.Kv.Key
				delete(m.Nodes,string(k))
			default:
				fmt.Printf("[%s] dir:%q, value:%q\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
			}
		}
	}
}

func (m *EtcdMaster)GetKeyInfo(key []byte) (*KeyInfo,error) {
	keyInfo := &KeyInfo{}
	err := keyInfo.SetRegisterKey(string(key))
	return keyInfo,err
}

func (m *EtcdMaster)GetValueInfo(value []byte) (*ValueInfo,error) {
	valueInfo := &ValueInfo{}
	err := valueInfo.DecodeValue(value)
	return valueInfo,err
}

//添加节点
func (m *EtcdMaster) addNode(key string,keyInfo *KeyInfo,valueInfo *ValueInfo){
	node := &EtcdNode{
		State:   valueInfo.Status == 0,
		KeyInfo:keyInfo,
		ValueInfo:valueInfo,
	}
	m.Nodes[key] = node
}

func (m *EtcdMaster)GetAllNodes() []EtcdNode  {
	var temp []EtcdNode
	for _,v := range m.Nodes{
		if v != nil{
			temp = append(temp,*v)
		}
	}
	return temp
}

func (m *EtcdMaster) GetNodeRandom()(EtcdNode,bool){
	count := len(m.Nodes)
	if count == 0{
		return EtcdNode{},false
	}
	idx := rand.Intn(count)
	for _,v := range m.Nodes{
		if idx == 0{
			return *v,true
		}
		idx = idx-1
	}
	return EtcdNode{},false
}