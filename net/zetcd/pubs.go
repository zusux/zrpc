package zetcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/sd"
)

type Pubs map[string]*Etcd

func (p Pubs) Register()  {
	for _,v :=range p{
		v.Register()
	}
}
func (p Pubs) UnRegister()  {
	for _,v :=range p{
		v.UnRegister()
	}
}
func (p Pubs) RequestLocal(ctx context.Context,kind string ,req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	etcd, ok := p[ kind ]
	if ok {
		return etcd.RequestLocal(ctx,req,reqFactory)
	} else {
		return nil, errors.New(fmt.Sprintf("没有注册本地类别: %s",kind))
	}
}
func (p Pubs) GrpcRequestRemote(ctx context.Context, serverName string, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	etcd, ok := p[ "grpc" ] /*如果确定是真实的,则存在,否则不存在 */
	if ok {
		return etcd.RequestRemote(ctx,serverName,"grpc",req,reqFactory)
	} else {
		//没有现成的则新建一个client
		return NewEtcdForClient().RequestRemote(ctx,serverName,"grpc",req,reqFactory)
	}
}

func (p Pubs) HttpRequestRemote(ctx context.Context, serverName string, req interface{}, reqFactory sd.Factory) (response interface{}, err error) {
	etcd, ok := p[ "http" ] /*如果确定是真实的,则存在,否则不存在 */
	if ok {
		return etcd.RequestRemote(ctx,serverName,"http",req,reqFactory)
	} else {
		//没有现成的则新建一个client
		return NewEtcdForClient().RequestRemote(ctx,serverName,"http",req,reqFactory)
	}
}
