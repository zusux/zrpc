package zrpc

import (
	"fmt"
	"github.com/zusux/zrpc/zerr"
)

type Error struct {
	Status  zerr.TLSType
	Code    string
	Message string
}

func NewError(status zerr.TLSType, message string) *Error  {
	return &Error{
		Status: status,
		Code: status.String(),
		Message: message,
	}
}

func (e *Error) Error() string  {
	return e.Message
}
//time="2021-09-24T20:52:10+08:00" level=info msg="初始化Etcd配置成功" init=etcd
func (e *Error) String() string  {
	return fmt.Sprintf(`status=%d code=%s message=%s`,e.Status,e.Code,e.Message)
}

func (e *Error) GetStatus() zerr.TLSType {
	return e.Status
}

func (e *Error) GetCode() string {
	return e.Code
}