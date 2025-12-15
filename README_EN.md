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

[![LATEST_VERSION](https://img.shields.io/github/v/tag/keepchen/go-sail?color=ff0000&amp;&amp;logo=go&amp;label=LATEST_VERSION)](https://github.com/keepchen/go-sail/tags)  

English | [简体中文](./README.md)

## Whats the go-sail?

**go-sail** is a lightweight progressive web framework implemented using Go language. **It is not the product of reinventing the wheel**, but stands on the shoulders of giants and integrates existing excellent components to help users build stable and reliable services in the simplest way.
As its name suggests, you can regard it as the beginning of your own journey in the golang ecosystem. go-sail will help you start lightly and set sail.

## How to use
> go version >= 1.23

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


## Examples
### Configuration
```go
parseFn := func(content []byte, viaWatch bool){
    fmt.Println("config content: ", string(content))
    if viaWatch {
        //reload config...
    }
}
etcdConf := etcd.Conf{
	Endpoints: []string{""},
	Username: "",
	Password: "",
}
key := "go-sail.config.yaml"

sail.Config(true, parseFn).ViaEtcd(etcdConf, key).Parse(parseFn)
```
### Log trace
```go
func UserRegisterSvc(c *gin.Context) {
  ...
  sail.LogTrace(c).Warn("log something...")
  ...
}
```
### JWT authentication
- Issue token
```go
func UserLoginSvc(c *gin.Context) {
  ...
  uid := "user-1000"
  exp := time.Now().Add(time.Hour * 24).Unix()
  otherFields := map[string]interface{}{
      "nickname": "go-sail",
      "avatar": "https://go-sail.dev/assets/avatar/1.png",
      ...
  }
  ok, token, err := sail.JWT().MakeToken(uid, exp, otherFields)
  ...
}
```
- Authentication
```go
func UserInfoSvc(c *gin.Context) {
  ...
  ok, claims, err := sail.JWT().ValidToken(token)
  ...
}
```
### Components
#### Responder
```go
func UserInfoSvc(c *gin.Context) {
  sail.Response(c).Wrap(constants.ErrNone, resp).Send()
}
```

#### Database
- Read / Write
```go
func UserInfoSvc(c *gin.Context) {
  uid := "user-1000"
  var user models.User
  //READ: query user info
  sail.GetDBR().Where("uid = ?", uid).First(&user)
  ...
  //WRITE: update user info
  sail.GetDBW().Model(&models.User{}).
      Where("uid = ?", uid).
      Updates(map[string]interface{}{
          "avatar": "https://go-sail.dev/assets/avatar/2.png"
      })
}
```
- Transaction
```go
func UserInfoSvc(c *gin.Context) {
  uid := "user-1000"
  err := sail.GetDBW().Transaction(func(tx *gorm.DB){
      e1 := tx.Model(&models.User{}).
              Where("uid = ?", uid).
              Updates(map[string]interface{}{
                  "avatar": "https://go-sail.dev/assets/avatar/2.png"
              }).Error
      if e1 != nil {
          return e1
      }
      e2 := tx.Create(&models.UserLoginHistory{
                Uid: uid,
                ...
              }).Error
      return e2
  })
}
```
#### Redis
```go
func UserInfoSvc(c *gin.Context) {
  ...
  sail.GetRedis().Set(ctx, "go-sail:userInfo", "user-1000", time.Hour*24).Result()
  ...
}
```
### Task schedule
- Interval
```go
func TodoSomething() {
  fn := func() { ... }
  sail.Schedule("todoSomething", fn).Daily()
}
```
- Linux Crontab style
```go
func TodoSomething() {
  fn := func() { ... }
  sail.Schedule("todoSomething", fn).RunAt("*/5 * * * *")
}
```
- Race detection
```go
func TodoSomething() {
  fn := func() { ... }
  sail.Schedule("todoSomething", fn).Withoutoverlapping().RunAt("*/5 * * * *")
}
```
### Distributed lock
```go
func UpdateUserBalance() {
  if !sail.RedisLocker().TryLock(key) {
      return false
  }
  defer sail.RedisLocker().Unlock(key)
  ...
}
```

## Documentation
[https://go-sail.dev](https://go-sail.dev)  

## Live demo  
[https://nav.go-sail.dev](https://nav.go-sail.dev)

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
    - Valkey
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
- [x] Configuration
  - File
  - Etcd
  - Nacos

#### Other Plugins
[README.md](plugins/README.md)

## Benchmark
```shell
ulimit -n 65535 && sh run_benchmark.sh
```  
Test results (real HTTP requests)  
```text
goos: darwin
goarch: amd64
pkg: github.com/keepchen/go-sail/v3
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGoSailParallel-12    88252    12898 ns/op    8860 B/op    92 allocs/op
BenchmarkGinParallel-12       96548    11722 ns/op    7187 B/op    82 allocs/op
PASS
ok    github.com/keepchen/go-sail/v3  3.663s
```  
![benchmark-result.png](static/benchmark-result.png)

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
  <tr style="height:200px">
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://stardots.io?ref=go-sail" target="_blank"><img src="static/usecases/stardots-logo.png" alt="stardots.io" width="200" title="https://stardots.io"/></a>
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://t.me/PiggyPiggyofficialbot" target="_blank"><img src="static/usecases/piggy-logo.jpg" alt="Piggy (telegram mini-game)" width="200" title="https://t.me/PiggyPiggyofficialbot"/></a>
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/miniprogram-hpp.png" alt="生活好评助手-小程序" width="200" />
    </td>
  </tr>
  <tr style="height:200px">
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/wingoal-metaland.png" alt="WinGoal" width="200" />
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <img src="static/usecases/pikaster-metaland.png" alt="Pikaster" width="200" />
    </td>
    <td style="border: 1px solid black; padding: 8px;">
      <a href="https://fantagoal.io" target="_blank"><img src="static/usecases/fantaGoal-logo.png" alt="FantaGoal" width="200" title="https://fantagoal.io"/></a>
    </td>
  </tr>

</table>

## Sponsors
[![Powered by DartNode](static/sponsors/DartNode_Brand_Full/black_color_full.png)](https://dartnode.com "Powered by DartNode - Free VPS for Open Source")

## Star History  
[![Star History Chart](https://api.star-history.com/svg?repos=keepchen/go-sail&type=date&legend=top-left)](https://www.star-history.com/#keepchen/go-sail&type=date&legend=top-left)  
