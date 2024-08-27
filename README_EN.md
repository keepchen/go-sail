<div align="center">
    <h1><img src="static/sailboat-solid-colorful.svg" alt="sailboat-solid" title="sailboat-solid" width="300" /></h1>
</div> 

[![Go](https://github.com/keepchen/go-sail/actions/workflows/go.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/go.yml)  [![CodeQL](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml)  [![Go Report Card](https://goreportcard.com/badge/github.com/keepchen/go-sail/v3)](https://goreportcard.com/report/github.com/keepchen/go-sail/v3)  

English | [简体中文](./README.md)

## Whats the go-sail？

**go-sail** is a lightweight progressive web framework implemented using Go language. **It is not the product of reinventing the wheel**, but stands on the shoulders of giants and integrates existing excellent components to help users build stable and reliable services in the simplest way.
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
    conf = &config.Config{}
    registerRoutes = func(ginEngine *gin.Engine) {
        ginEngine.GET("/hello", func(c *gin.Context){
            c.String(http.StatusOK, "%s", "hello, world!")
        })
    }
)

func main() {
    sail.WakeupHttp("go-sail", conf).Hook(registerRoutes, nil, nil).Launch()
}
```  
Console screenshot after launched like this:  

<img src="static/launch.png" alt="launch.png" title="launch.png" width="600" />  

## Documentation
[Docs](https://go-sail.keepchen.com)

## Features  
- [x] HTTP Responder
    - Uniform Response Fields
    - Managing HTTP status codes  
    - Management business code  
- [x] Highly used component library 
    - database
    - email
    - jwt
    - kafka
    - logger
    - nacos
    - etcd
    - nats
    - redis
- [x] Service Registration and Discovery  
    - Nacos
    - Etcd
- [x] Frequently used tools  
    - Encryption and Decryption  
    - File
    - IP
    - String
    - Random number  
    - Date and Time  
    - ...
- [x] Log collection and export
    - Local files
    - Exporter  
- [x] Scheduled Tasks  
    - Cancellable  
    - Disposable  
    - Periodic  
    - Linux Crontab style  
    - Race Detection  
- [x] Call chain log tracking  
    - Passing through the request context  
- [x] Multi-language error codes
    - Dynamic injection  
- [x] Distributed lock based on Redis  
    - Standalone mode  
    - Cluster mode
- [x] API Documentation
    - Redocly
    - Swagger  

#### Other Plugins
[README.md](plugins/README.md)  

## Big Thanks
Thank you to everyone who provided valuable suggestions and comments during the experience and use, as well as provided other kinds of help!
- Configuration modular optimization [@fujilin](https://github.com/fujilin)
- Responder syntax sugar enhancement optimization [@lichuanzhang](https://github.com/lichuanzhang)
- Logo beautification [@ShuaiRen34](https://twitter.com/ShuaiRen34)

## Other
- PR is welcome: [pull request](https://github.com/keepchen/go-sail/compare)
- Issue is welcome: [issue](https://github.com/keepchen/go-sail/issues/new/choose)
- Thank you for your star if you like this project :)  

## Use cases
<img src="static/usecases/stardots-logo.png" alt="stardots.ink" width="300" title="https://stardots.ink"/>
<br/><br/>
<img src="static/usecases/piggy-logo.jpg" alt="Piggy (telegram mini-game)" width="300" title="https://t.me/PiggyPiggyofficialbot"/>
<br/><br/>
<img src="static/usecases/fantaGoal-logo.png" alt="FantaGoal" width="300" title="https://fantagoal.io"/>
<br/><br/>
<img src="static/usecases/pikaster-metaland.png" alt="Pikaster" width="300" />
<br/><br/>
<img src="static/usecases/wingoal-metaland.png" alt="WinGoal" width="300" />
<br/><br/>
<img src="static/usecases/miniprogram-hpp.png" alt="生活好评助手-小程序" width="180" />