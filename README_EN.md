<div align="center">
    <h1><img src="static/sailboat-solid-colorful.svg" alt="sailboat-solid" title="sailboat-solid" width="300" /></h1>
</div> 

[![Go](https://github.com/keepchen/go-sail/actions/workflows/go.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/go.yml)  [![CodeQL](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml)  [![Go Report Card](https://goreportcard.com/badge/github.com/keepchen/go-sail/v3)](https://goreportcard.com/report/github.com/keepchen/go-sail/v3)  

English | [简体中文](./README.md)

## Whats the go-sail？

**go-sail** is a lightweight progressive web framework implemented using Go language. It is not the product of reinventing the wheel, but stands on the shoulders of giants and integrates existing excellent components to help users build stable and reliable services in the simplest way.
As its name suggests, you can regard it as the beginning of your own journey in the golang ecosystem. go-sail will help you start lightly and set sail.

## How to use
> go version >= 1.19  

> go get -u github.com/keepchen/go-sail/v3

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/keepchen/go-sail/v3/sail"
    "github.com/keepchen/go-sail/v3/sail/config"
)

var (
    conf = &config.Config{
        LoggerConf: logger.Conf{
            Filename: "logs/running.log",
        },
        HttpServer: config.HttpServerConf{
            Debug: true,
            Addr:  ":8000",
        },
    }
    registerRoutes = func(ginEngine *gin.Engine) {
        ginEngine.GET("/hello", func(c *gin.Context){
            c.String(http.StatusOK, "%s", "hello, world!")
        })
    }
)

sail.WakeupHttp("go-sail", conf).
    Hook(registerRoutes, nil, nil).
    Launch()
```  
Console screenshot after launched like this:  

<img src="static/launch.png" alt="launch.png" title="launch.png" width="600" />  

## Documentation
[Docs](https://blog.keepchen.com/a/go-sail.html?from=github)

## Features
- Get components
> go-sail启动时，会根据配置文件启动相应的应用组件，可使用`sail`关键字统一获取
```go
import (
    "github.com/keepchen/go-sail/v3/sail"
)

//获取日志组件
sail.GetLogger()

//获取数据库连接（读、写实例）
sail.GetDB()

//获取redis连接(单例模式)
sail.GetRedis()

//获取redis连接(cluster模式)
sail.GetRedisCluster()

//获取nats连接
sail.GetNats()

//获取kafka完整连接实例
sail.GetKafkaInstance()

//获取etcd连接实例
sail.GetEtcdInstance()
```  
PR is welcome👏🏻👏🏻

- Response
```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/keepchen/go-sail/v3/constants"
    "github.com/keepchen/go-sail/v3/sail"
)

//handler
func SayHello(c *gin.Context) {
    sail.Response(c).Builder(constants.ErrNone, nil, "OK").Send()
}
```  

- Response (entity)
```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/keepchen/go-sail/v3/constants"
    "github.com/keepchen/go-sail/v3/http/pojo/dto"
    "github.com/keepchen/go-sail/v3/sail"
)

type UserInfo struct {
    dto.Base
    Data struct {
        Nickname string `json:"nickname" validate:"required" format:"string"` //nickname
        Age int `json:"int" validate:"required" format:"number"` //age
    } `json:"data" validate:"required"` //body data
}

// implement dto.IResponse interface
func (v UserInfo) GetData() interface{} {
    return v.Data
}

//handler
func GetUserInfo(c *gin.Context) {
    var resp UserInfo
    resp.Data.Nickname = "go-sail"
    resp.Data.Age = 18
	
    sail.Response(c).Builder(constants.ErrNone, resp).Send()
}
```

#### Other Plugins
[README.md](plugins/README.md)

## Use cases
<img src="static/usecases/pikaster-metaland.png" alt="Pikaster" width="300" />
<img src="static/usecases/wingoal-metaland.png" alt="WinGoal" width="300" />
<img src="static/usecases/miniprogram-hpp.png" alt="生活好评助手-小程序" width="180" />