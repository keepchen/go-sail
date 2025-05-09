<div align="center">
    <h1><img src="static/sailboat-solid-colorful.svg" alt="sailboat-solid" title="sailboat-solid" width="300" /></h1>
</div> 

[![Go](https://github.com/keepchen/go-sail/actions/workflows/go.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/go.yml)
[![CodeQL](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml/badge.svg)](https://github.com/keepchen/go-sail/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/keepchen/go-sail/v3)](https://goreportcard.com/report/github.com/keepchen/go-sail/v3)
[![codecov](https://codecov.io/github/keepchen/go-sail/graph/badge.svg?token=UNLOORRJHA)](https://codecov.io/github/keepchen/go-sail)
[![Dependabot Status](https://img.shields.io/badge/dependabot-active-brightgreen?logo=dependabot)](https://github.com/keepchen/go-sail/security/dependabot)
[![Snyk Security](https://img.shields.io/badge/Snyk-Secure-blueviolet?logo=snyk)](https://snyk.io/test/github/keepchen/go-sail)
[![LICENSE: MIT](https://img.shields.io/github/license/keepchen/go-sail.svg?style=flat)](LICENSE)    

English | [简体中文](./README.md)

## Whats the go-sail？

**go-sail** is a lightweight progressive web framework implemented using Go language. **It is not the product of reinventing the wheel**, but stands on the shoulders of giants and integrates existing excellent components to help users build stable and reliable services in the simplest way.
As its name suggests, you can regard it as the beginning of your own journey in the golang ecosystem. go-sail will help you start lightly and set sail.

## How to use
> go version >= 1.20  

> go get -u github.com/keepchen/go-sail/v3

```go  
import (
    "net/http"
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
[Docs](https://go-sail.dev)

## Features  
- [x] HTTP Responder
    - Uniform Response Fields
    - Managing HTTP status codes  
    - Management business code  
- [x] Components  
    - Database
    - Email
    - Jwt
    - Kafka
    - Logger
    - Nacos
    - Etcd
    - Nats
    - Redis
- [x] Service Registration and Discovery  
    - Nacos
    - Etcd
- [x] Toolkit  
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
      - Redis  
      - Kafka  
      - Nats  
- [x] Scheduled Tasks  
    - Cancellable  
    - Disposable  
    - Periodic  
    - Linux Crontab style  
    - Race Detection  
- [x] Telemetry and Observability
  - Call chain tracing
  - Prometheus
  - Pprof
  - Log Exporter  
- Performance monitor  
  - Prometheus  
  - Pprof  
- [x] API error codes
    - Dynamic injection  
    - Internationalization  
- [x] Distributed lock based on Redis  
  - Blocking  
  - None-Blocking  
- [x] API Documentation
    - Redocly
    - Swagger  

#### Other Plugins
[README.md](plugins/README.md)  

## Big Thanks
Thank you to everyone who provided valuable suggestions and comments during the experience and use, as well as provided other kinds of help!
- Configuration modular optimization proposal [@fujilin](https://github.com/fujilin)
- Responder syntax sugar enhancement optimization proposal [@lichuanzhang](https://github.com/lichuanzhang)
- Logo beautification [@ShuaiRen34](https://twitter.com/ShuaiRen34)

## Other
- PR is welcome: [pull request](https://github.com/keepchen/go-sail/compare)
- Issue is welcome: [issue](https://github.com/keepchen/go-sail/issues/new/choose)
- Thank you for your star if you like this project :)  

## Use cases  
<table style="text-align: center">
  <tr>
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://stardots.io?ref=go-sail" target="_blank"><img src="static/usecases/stardots-logo.png" alt="stardots.io" width="200" title="https://stardots.io"/></a>
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/pikaster-metaland.png" alt="Pikaster" width="200" />
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://fantagoal.io" target="_blank"><img src="static/usecases/fantaGoal-logo.png" alt="FantaGoal" width="200" title="https://fantagoal.io"/></a>
    </td>
  </tr>
  <tr>
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/wingoal-metaland.png" alt="WinGoal" width="200" />
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://t.me/PiggyPiggyofficialbot" target="_blank"><img src="static/usecases/piggy-logo.jpg" alt="Piggy (telegram mini-game)" width="200" title="https://t.me/PiggyPiggyofficialbot"/></a>
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/miniprogram-hpp.png" alt="生活好评助手-小程序" width="200" />
    </td>
  </tr>
</table>  

## Star History  
![](https://starchart.cc/keepchen/go-sail.svg)