package zetcd

type Etcd struct {
	Hosts         []string
	DialTimeout   int64
	DialKeepalive int64
}

//SetDialTimeout 设置超时时间
func (z *Etcd) getDialTimeout() int64 {
	if z.DialTimeout == 0{
		z.DialTimeout = 500
	}
	return z.DialTimeout
}

//getDialKeepalive 设置超时时间
func (z *Etcd) getDialKeepalive() int64 {
	if z.DialKeepalive == 0{
		z.DialKeepalive = 10
	}
	return z.DialKeepalive
}

//SetDialTimeout 设置超时时间
func (z *Etcd) SetDialTimeout(dialTimeout int64) *Etcd {
	z.DialTimeout = dialTimeout
	return z
}
//SetDialKeepAlive 设置keepalive时间
func (z *Etcd) SetDialKeepAlive(dialKeepAlive int64) *Etcd {
	z.DialKeepalive = dialKeepAlive
	return z
}