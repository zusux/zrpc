package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/zusux/zrpc"
	"github.com/zusux/zrpc/example/access"
	"google.golang.org/grpc"
	"io"
	"time"
)

func main()  {
	grpcReq()
}

func etcd()  {
	conf := zrpc.GetConfig()
	zrpc.Log().Info(fmt.Sprintf("%+v",conf))
	conf.Reg.Register()
	defer conf.Reg.UnRegister()
	for i:=0;i<10;i++{
		time.Sleep(time.Second)
	}
	fmt.Println("close")
}

func db()  {
	db := zrpc.GetDb()
	type Book struct {
		Id   int
		Name string
	}
	var book Book
	db.First(&book)
	zrpc.Log().Info(fmt.Sprintf("book: %+v", book))
}

func log() {
	log := zrpc.Log()
	log.Info("aaa")
	log.Error("fff")
	log.Warning("www")
	log.Info(zrpc.K.String("server.name"))
}

//etcd grpc
func grpcReq() {
	zrpc.GetEtcd().Register()
	doReq()
	zrpc.GetEtcd().UnRegister()
}

func doReq() {
	var req interface{}
	zrpc.GetEtcd().GrpcRequestRemote(context.Background(), "github.com/zusux/go-doc/ibook-service",req, reqFactory)
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