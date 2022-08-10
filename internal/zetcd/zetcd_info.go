package zetcd

import (
	"errors"
	"fmt"
	"github.com/zusux/zrpc/internal/proto/pb_value"
	"google.golang.org/protobuf/proto"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type KeyInfo struct {
	Cluster string //集群
	Name    string //服务名称
	Kind    string   //grpc,http
}

type ValueInfo struct {
	Kind    string   //grpc,http
	Ip      string //服务注册ip
	Port    uint32    //服务注册端口
	Status  STATUS //0 初始化, 1 正常 , 2 挂起, 3关闭
	RequestFlow uint32
	UpdatedAt int64
}

func NewKeyInfo(cluster,name,kind string) *KeyInfo {
	return &KeyInfo{
		Cluster: cluster,
		Name:    name,
		Kind:    kind,
	}
}

func NewValueInfo(kind,ip string, port uint32 ,status STATUS, requestFlow uint32,updatedAt int64) *ValueInfo {
	return &ValueInfo{
		Kind:    kind,
		Ip:    ip,
		Port:    port,
		Status:    status,
		RequestFlow:    requestFlow,
		UpdatedAt:    updatedAt,
	}
}


func (i *KeyInfo) GetRegisterKey() string {
	return fmt.Sprintf("%s/%s:%s", i.Cluster,i.Name, i.Kind)
}

func (i *KeyInfo) GetNameKind() string {
	return fmt.Sprintf("%s:%s",i.Name, i.Kind)
}

// SetRegisterKey 设置注册地址
func (i *KeyInfo) SetRegisterKey(address string) error {
	lindex := strings.LastIndex(address,":")
	if lindex == -1 {
		return errors.New("[zetcd] not match format name:kind")
	}
	k := strings.TrimLeft(address[lindex:],":")
	i.setKind(k)
	findex := strings.Index(address,"/")
	if findex == -1 || findex>= lindex {
		return errors.New("[zetcd] not match format cluster/name")
	}
	cluster := address[0:findex]
	i.setCluster(cluster)
	name := address[findex+1:lindex]
	i.setName(name)
	return nil
}

// setCluster 设置集群
func (i *KeyInfo) setCluster(cluster string) {
	i.Cluster = cluster
}
// setName 设置名称
func (i *KeyInfo) setName(name string) {
	i.Name = name
}
// setKind 设置类型
func (i *KeyInfo) setKind(kind string) {
	i.Kind = kind
}

func (i *ValueInfo) getRegisterAddress() string {
	return fmt.Sprintf("%s:%d", i.Ip, i.Port)
}

func (i *ValueInfo) GetValidAddress() (string,bool) {
	str := fmt.Sprintf("%s:%d", i.Ip, i.Port)
	result,_  :=regexp.MatchString(`\d{1,3}\.\d{1,3}\.\d{1,3}.\d{1,3}\:\d+`,str)
	return str,result
}

// SetRegisterAddress 设置服务地址 ip:port
func (i *ValueInfo) SetRegisterAddress(address string) error {
	arr := strings.Split(address, ":")
	if len(arr) != 2 {
		return errors.New("[zetcd] not match format ip:port")
	}
	ip := net.ParseIP(arr[0])
	if ip == nil {
		return errors.New("[zetcd] not match ip format")
	}
	i.setIp(ip.String())
	port, err := strconv.Atoi(arr[1])
	if err != nil {
		return errors.New(fmt.Sprintf("[zetcd] parse port error:%s",err.Error()))
	}
	i.setPort(uint32(port))
	return nil
}
// EncodeValue 编码 proto消息
func (i *ValueInfo) EncodeValue() ([]byte,error) {
	status := uint32(i.Status)
	updatedAt := uint32(i.UpdatedAt)
	value := pb_value.InfoValue{
		Kind: &i.Kind,
		Ip: &i.Ip,
		Port: &i.Port,
		Status: &status,
		RequestFlow: &i.RequestFlow,
		UpdatedAt: &updatedAt,
	}
	// 将protocol消息序列化成二进制
	valueByte, err := proto.Marshal(&value)
	if err != nil {
		return nil,errors.New(fmt.Sprintf("[zetcd] decode value error:%s\n",err.Error()))
	}
	fmt.Println(valueByte)
	return valueByte,nil
}
// DecodeValue 解码 proto消息
func (i *ValueInfo) DecodeValue(infoByte []byte) error  {
	// 将二进制消息反序列化
	newValue := pb_value.InfoValue{}
	err := proto.Unmarshal(infoByte, &newValue)
	if err != nil {
		return errors.New(fmt.Sprintf("[zetcd] uncode value error:%s\n",err.Error()))
	}
	i.Ip = *newValue.Ip
	i.Port = *newValue.Port
	i.Kind = *newValue.Kind
	i.RequestFlow = *newValue.RequestFlow
	i.Status = STATUS(*newValue.Status)
	i.UpdatedAt = int64(*newValue.UpdatedAt)
	fmt.Printf("%+v\n", newValue)
	fmt.Println(newValue)
	return nil
}

// setIp 设置ip
func (i *ValueInfo) setIp(ip string) {
	i.Ip = ip
}

// setPort 设置端口
func (i *ValueInfo) setPort(port uint32) {
	i.Port = port
}


