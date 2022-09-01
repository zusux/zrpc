package internal

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Host string
	Port int
	Auth string
	Db   int
}

func NewRedis(host string, port int, auth string, db int) *Redis {
	return &Redis{
		Host: host,
		Port: port,
		Auth: auth,
		Db:   db,
	}
}

func (r *Redis) Client() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.getAddr(),
		Password: r.Auth, // no password set
		DB:       r.Db,   // use default DB
	})
	return rdb
}

func (r *Redis) getAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
