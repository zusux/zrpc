package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/zusux/zrpc"
	"github.com/zusux/zrpc/env"
	"github.com/zusux/zrpc/example/access"
	"github.com/zusux/zrpc/net/zetcd"
	"google.golang.org/grpc"
	"io"
	"os"
)

//log
func main() {
	main_config()
}

func main_config() {
	zrpc.Init()
	var servers map[string]zetcd.Server
	// Quick unmarshal.
	env.K.Unmarshal("zrpc", &servers)
	fmt.Println(servers)

}

func main_log() {
	os.Setenv("ENV_FILE", "./../env/env.toml") //env配置文件实际位置
	zrpc.Init()
	zrpc.GetLog().Info("hello")
}

//gorm
func main_db() {
	os.Setenv("ENV_FILE", "./../env/env.toml") //env配置文件实际位置
	zrpc.Init()
	gorm := zrpc.GetDb()
	type Book struct {
		Id   int64
		Name string
	}
	var book Book
	gorm.Table("book").First(&book)
	fmt.Println(book)
}

//etcd grpc
func main_grpc() {
	os.Setenv("ENV_FILE", "./../env/env.toml") //env配置文件实际位置
	zrpc.Init()
	zrpc.GetEtcd().Register()
	doReq()
	zrpc.GetEtcd().UnRegister()
}

func doReq() {
	var req interface{}
	zrpc.GetEtcd().RequestLocal(context.Background(), "grpc",req, reqFactory)
}

//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	//fmt.Println("instanceAddr",instanceAddr)
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//fmt.Println("请求服务: ", instanceAddr)
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println("dial error", err)
			return nil, err
		}
		defer conn.Close()
		bookClient := access.NewIBookClient(conn)
		_, err = bookClient.BookInfo(context.Background(), &access.BookInfoReq{BookId: 2})
		if err != nil {
			fmt.Println("BookInfoClient", err)
			return nil, err
		}
		//fmt.Println("获取书籍详情")
		//fmt.Println("bookId: 2", " => ", "bookName:", bi.Data.Name)

		bl, err := bookClient.BookList(context.Background(), &access.BookListReq{Page: 1, Limit: 9})
		if err != nil {
			fmt.Println("bookListClient", err)
			return nil, err
		}
		fmt.Println("获取书籍列表")
		if bl.Data != nil {
			for _, b := range bl.Data {
				fmt.Println("bookId:", b.Id, " => ", "bookName:", b.Name)
			}
		} else {
			fmt.Println("bl.Message", bl.Message)
		}
		return nil, nil
	}, nil, nil
}
