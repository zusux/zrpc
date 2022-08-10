package internal

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error  {
	return &Error{
		Code: code,
		Message: message,
	}
}

func (e *Error) Error() string  {
	return e.Message
}
//time="2021-09-24T20:52:10+08:00" level=info msg="初始化Etcd配置成功" init=etcd
func (e *Error) String() string  {
	return fmt.Sprintf(`code:%d message:%s`,e.Code,e.Message)
}

func (e *Error) GetCode() int {
	return e.Code
}