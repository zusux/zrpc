package internal

import (
	"github.com/zusux/zrpc/internal/util"
	"log"
)



type Server struct {
	Cluster string
	Name    string
	Etcd    struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepalive int64
	}
	Ip     string
	Retry  bool
	Listen map[string]Pub
}

func GetServer() (*Server, error) {
	svr := Server{}
	// Quick unmarshal.
	err := K.Unmarshal("server", &svr)
	if err != nil {
		log.Printf("[zrpc][error] unmarshal config err: %s",err)
		return &svr, err
	}
	return &svr, nil
}

func (s *Server) Reg() *Reg {

	if s.Ip == "" {
		ip, err := util.GetLocalIP()
		if err != nil {
			log.Printf("[zrpc] reg get local ip err: %s",err)
			return nil
		}
		s.Ip = ip.String()
	}
	if len(s.Etcd.Hosts) > 0 {
		reg := NewReg(s.Cluster, s.Name, s.Ip, s.Etcd.DialTimeout, s.Etcd.DialKeepalive, s.Retry, s.Etcd.Hosts, s.Listen)
		return reg
	}
	return nil
}
