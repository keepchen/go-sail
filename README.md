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

简体中文 | [English](./README_EN.md)

## go-sail是什么？  

**go-sail**是一个轻量的渐进式Web框架，使用Go语言实现。它并**不是重复造轮子的产物**，而是站在巨人的肩膀上，整合现有的优秀组件，旨在帮助使用者以最简单的方式构建稳定可靠的服务。  
正如它的名字一般，你可以把它视作自己在golang生态的一个开始。go-sail将助力你从轻出发，扬帆起航。  

## 如何使用  
> 推荐go version >= 1.19  

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
当你看到终端如下图所示内容就表示服务启动成功了：  

<img src="static/launch.png" alt="launch.png" title="launch.png" width="600" />  

## 文档  
[文档传送门](https://go-sail.keepchen.com)

## 功能特性  
- [x] HTTP响应器  
  - 统一响应字段  
  - 管理HTTP状态码  
  - 管理业务码  
- [x] 组件库  
  - Database  
  - Email  
  - Jwt  
  - Kafka  
  - Logger  
  - Nacos  
  - Etcd  
  - Nats  
  - Redis  
- [x] 服务注册与发现  
  - Nacos  
  - Etcd  
- [x] 工具类  
  - 加解密  
  - 文件  
  - ip  
  - 字符串  
  - 随机数  
  - 日期时间  
  - ...
- [x] 日志收集与导出  
  - 本地文件  
  - 导出器  
    - Redis  
    - Kafka  
    - Nats  
- [x] 计划任务  
  - 可取消的  
  - 一次性的  
  - 周期性的  
  - Linux Crontab风格的  
  - 竞态检测  
- [x] 遥测与可观测性  
  - 调用链追踪
  - Prometheus  
  - Pprof 
  - 日志导出器    
- [x] 接口错误码  
  - 动态注入  
  - 国际化  
- [x] 基于Redis的分布式锁  
  - 阻塞式  
  - 非阻塞式  
- [x] 接口文档  
  - Redocly  
  - Swagger

#### 其他插件  
[README.md](plugins/README.md)  

## 大感谢  
感谢在体验、使用过程中提出宝贵建议和意见以及提供过其他各种帮助的各位小伙伴！  
- 配置模块化优化建议 [@fujilin](https://github.com/fujilin)  
- 响应器语法糖增强优化建议 [@lichuanzhang](https://github.com/lichuanzhang)  
- Logo美化 [@ShuaiRen34](https://twitter.com/ShuaiRen34)  

## 其他  
- 欢迎大家提PR: [pull request](https://github.com/keepchen/go-sail/compare)  
- 欢迎大家提出自己的想法: [issue](https://github.com/keepchen/go-sail/issues/new/choose)  
- 感谢你的star如果你喜欢这个项目的话 :)  

## 使用案例  
<a href="https://stardots.io?ref=go-sail" target="_blank"><img src="static/usecases/stardots-logo.png" alt="stardots.io" width="300" title="https://stardots.io"/></a>
<br/><br/>
<a href="https://t.me/PiggyPiggyofficialbot" target="_blank"><img src="static/usecases/piggy-logo.jpg" alt="Piggy (telegram mini-game)" width="300" title="https://t.me/PiggyPiggyofficialbot"/></a>
<br/><br/>
<a href="https://fantagoal.io" target="_blank"><img src="static/usecases/fantaGoal-logo.png" alt="FantaGoal" width="300" title="https://fantagoal.io"/></a>
<br/><br/>
<img src="static/usecases/pikaster-metaland.png" alt="Pikaster" width="300" />
<br/><br/>
<img src="static/usecases/wingoal-metaland.png" alt="WinGoal" width="300" />
<br/><br/>
<img src="static/usecases/miniprogram-hpp.png" alt="生活好评助手-小程序" width="180" />  
 
![](https://starchart.cc/keepchen/go-sail.svg)

