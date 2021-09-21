package code

import "net/http"

func GetCode(err error,code int64) int64  {
	if err != nil{
		return code
	}else{
		return http.StatusOK
	}
}

func GetMessage(err error) string{
	if err != nil{
		return err.Error()
	}
	return ""
}