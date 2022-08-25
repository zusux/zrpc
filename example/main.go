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
	zrpc.Init()
	fmt.Println("inited ",zrpc.GetDb())
	//fmt.Println(os.Getenv("site_mode"))
	e := make(chan error,1000)
	for  {
		err :=db()
		if err !=nil{
			e<-err
		}
	}
	<-e
	//etcdRetry()
}

func etcd()  {
	conf := zrpc.GetConf()
	zrpc.Log().Info(fmt.Sprintf("%+v",conf))
	conf.Register()
	defer conf.UnRegister()
	for i:=0;i<10;i++{
		time.Sleep(time.Second)
		fmt.Println(i)
	}
	fmt.Println("close")
}

func etcdRetry()  {
	conf := zrpc.GetConf()
	zrpc.Log().Info(fmt.Sprintf("%+v",conf))
	conf.Register()
	defer conf.UnRegister()
	i:=0
	for {
		i++
		if i==5{
			for k,v := range conf.EtcdDis.MapRegister{
				zrpc.Log().Infof("Key %s, value: %v:%v",k, v.Value.Ip,v.Value.Port)
				v,err := v.GetInfo()
				zrpc.Log().Infof("info value: %+v, err:%v",v, err)
			}
		}

		if i==10{
			for k,v := range conf.EtcdDis.MapRegister{
				zrpc.Log().Infof("register Key %s, value: %v:%v",k, v.Value.Ip,v.Value.Port)
				//v.UpdateService(1,2)
				v,err := v.GetInfo()
				zrpc.Log().Infof("info value: %+v, err:%v",v, err)
			}
			i = 0
		}
		time.Sleep(time.Second)
		fmt.Println(i)

	}
	fmt.Println("close")
}

func db() error {
	db := zrpc.GetDb()
	type Book struct {
		Id   int
		Name string
	}
	var book Book
	err := db.First(&book).Error
	zrpc.Log().Info(fmt.Sprintf("book: %+v, err:%v", book, err))
	return err
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
	zrpc.GetConf().Register()
	doReq()
	zrpc.GetConf().UnRegister()
}

func doReq() {
	var req interface{}
	zrpc.GetConf().GrpcRequestRemote(context.Background(), "ibook-service",req, reqFactory)
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