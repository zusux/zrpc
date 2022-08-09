```shell
   go get github.com/zusux/zrpc
```

#### 初始化
```go
// 配置文件存放路径 默认在 执行目录/env/env.toml
//设置配置文件路径
os.SetEnv("env.config","env/env.toml")
//初始化
zrpc.Init()
```

#### log 获取
```go
log := zrpc.GetLog()
```


#### 配置项获取
```go
conf := zrpc.GetConf()
name := zrpc.K.String("server.name")
```   

#### gorm 获取
```go
   gorm  := zrpc.GetDb()
```

####grpc 服务注册
```go
   etcd := zrpc.GetEtcd()
   etcd.Register()
   etcd.GrpcRequest(ctx context.Context, req interface{}, reqFactory sd.Factory) 
   etcd.UnRegister()
```