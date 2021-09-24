package zerr

import (
	"fmt"
	"github.com/zusux/zrpc/code"
)

type Zerr struct {
	Status code.TLSType
	Code string
	Message string
}

func NewZErr(status code.TLSType, message string) *Zerr  {
	return &Zerr{
		Status: status,
		Code: status.String(),
		Message: message,
	}
}

func (z *Zerr) Error() string  {
	return z.Message
}
//time="2021-09-24T20:52:10+08:00" level=info msg="初始化Etcd配置成功" init=etcd
func (z *Zerr) String() string  {
	return fmt.Sprintf(`status=%d code=%s message=%s`,z.Status,z.Code,z.Message)
}

func (z *Zerr) GetStatus() code.TLSType {
	return z.Status
}

func (z *Zerr) GetCode() string {
	return z.Code
}