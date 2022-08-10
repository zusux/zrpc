package internal

type Redis struct{
	Host string
	Port int
	Auth string
}

func NewRedis(host string, port int, auth string) *Redis{
	return &Redis{
		Host: host,
		Port: port,
		Auth: auth,
	}
}