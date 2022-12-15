<div align="center">
    <h1><img src="./sailboat-solid.svg" alt="sailboat-solid" title="icon from font awesome" width="600" /></h1>
</div> 

简体中文 | [English](./README.md)

## go-sail是什么？  

**go-sail**可能是你见过的极少数小而美的渐进式golang工程目录，它并**不是一个框架**，而是站在巨人的肩膀上，将所需的组件适当整理，有机结合，
在保证生产可用的前提下做到尽可能精简整洁的**项目工程**，更多的是一种**想法**。使用子命令的方式，实现功能模块/服务的拆分，通过配置中心完成服务注册与自动发现。从单体架构向微服务时代迈进。  
`go-sail user`启动用户服务，`go-sail order`启动订单服务，……  
正如它的名字一般，你可以把它视作自己在golang生态的一个开始。go-sail将助力你从轻出发，扬帆起航。  

## 功能特点  

#### Http接口  
基于`gin-gonic/gin`http框架，具有轻量且高性能的特性，实现基本的路由注册、参数绑定、中间件挂载功能。

- 路由注册  

```go
r := gin.Default()
r.GET("/say-hello", handler.SayHello)
```  

- 参数绑定  

```go
var (
    form request.SayHello
    resp response.SayHello
)
if err := c.ShouldBind(&form); err != nil {
    api.New(c).Assemble(constants.ErrRequestParamsInvalid, nil).Send()
    return
}
```  

- 表单验证  
```go
var (
    form request.SayHello
    resp response.SayHello
)
if errorCode, err := form.Validator(); err != nil {
    api.New(c).Assemble(errorCode, nil, err.Error()).Send()
    return
}
```

- 统一返回  
```go
import "github.com/keepchen/go-sail/pkg/common/http/api"

//根据业务错误码自动设置http状态码
api.New(c).Assemble(constants.ErrNone, anyResponseData).Send() // <- 200
api.New(c).Assemble(constants.ErrRequestParamsInvalid, nil).Send() // <- 400
api.New(c).Assemble(constants.ErrInternalSeverError, nil).Send() // <- 500

//指定http状态码
api.New(c).Assemble(constants.ErrInternalSeverError, nil).SendWithCode(400) // <- 400

//自定义返回消息提示
api.New(c).Assemble(constants.ErrInternalSeverError, nil, "Whoops!Looks like something went wrong.").SendWithCode(400) // <- 400
```

- 基于路由中间件的日志请求参数打印、允许跨域、Prometheus指标记录  

```go
//全局打印请求载荷、放行跨域请求、写入Prometheus exporter
r.Use(mdlw.PrintRequestPayload(), mdlw.WithCors(allowHeaders), mdlw.PrometheusExporter())
```  

#### 日志组件  
go-sail基于`uber/zap`的日志类库和`natefinch/lumberjack`日志轮转类库，实现了按模块、分文件的日志记录功能，并支持配置文件启用基于redis list
的logstash导入方案。  

```go
//::初始化日志组件
logger.InitLoggerZap(config.GetGlobalConfig().Logger, "appName")

//::初始化日志组件（定义不同模块）
logger.InitLoggerZap(config.GetGlobalConfig().Logger, "appName", "api", "cron", "db")

//调用日志组件
logger.GetLogger().Info("hello~")

logger.GetLogger("api").Info("中间件:打印请求载荷", zap.Any("value", string(dump)))

logger.GetLogger("db").Error("数据库操作:CreateUserAndWallet:错误",
zap.Any("value", logger.MarshalInterfaceValue(userAndWallet)), zap.Errors("errors", []error{err}))
```  

#### 数据库组件  
go-sail基于`gorm.io/gorm`的数据库类库，实现了读写分离功能。得益于gorm丰富的driver支持，go-sail支持`mysql`、`sqlserver`、`postgresql`、`sqlite`、`clickhouse`数据库操作。  

```go
import "github.com/keepchen/go-sail/pkg/lib/db"

dbInstance := db.GetInstance()
dbR := dbInstance.R // <- 读实例
dbW := dbInstance.W // <- 写实例

err := dbR.Where(...).First(...).Error
err := dbW.Where(...).Updates(...).Error
```

#### 缓存组件  
go-sail基于`go-redis/redis`的redis类库，实现了对redis单实例和集群访问功能。  

```go
import "github.com/keepchen/go-sail/pkg/lib/redis"

redisInstance := redis.GetInstance()
redisInstance.Set(context.Background(), key, string(value), expired).Result()

redisClusterInstacne := redis.GetClusterInstance()
redisClusterInstacne.Set(context.Background(), key, string(value), expired).Result()
```

#### 配置中心  
go-sail基于`nacos-group/nacos-sdk-go`的配置中心类库，集成了配置热更、服务注册与发现功能。  

#### 文档工具  
go-sail基于`swaggo/swag`工具，实现了openapi文档生成功能。同时，go-sail提供了两种文档UI工具供你选择：  
1.基于`swaggo/gin-swagger`类库的Swagger UI  

<img src="./static/swagger-ui.png" alt="Swagger UI" />  

2.基于`Redocly/redoc`工具的Redoc UI  

<img src="./static/redoc-ui.png" alt="Redoc UI" />

#### 持续集成  
go-sail工程使用`harness/drone`CI/CD工具，实现对工程项目的自动化测试、集成与发布。参考[.drone.yml](./.drone.yml)文件配置。关于`drone`ci工具的部署和使用，如果你感兴趣，
请移步至 [GitLab+Drone使用体验](https://blog.keepchen.com/a/the-gitlab-drone-experience.html)。  

#### 构建与部署  
go-sail提供了`Dockerfile`docker镜像构建脚本，同时也提供了快速构建命令(shell命令)，帮助你快速方便的完成镜像构建。如需镜像仓库，
可以参考[keepchen/docker-compose](https://github.com/keepchen/docker-compose/tree/main/harbor)中关于harbor搭建的相关内容。
关于工程服务的快速启动，可以参考工程目录下的[docker-compose.yml](./docker-compose.yml)。

## 工程依赖  

#### 组件/类库  

- [spf13/cobra](https://github.com/spf13/cobra)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
- [gorm.io/gorm](https://github.com/go-gorm/gorm)
- [go.uber.org/zap](https://github.com/uber-go/zap)
- [go-redis/redis](https://github.com/go-redis/redis)
- [jinzhu/configor](https://github.com/jinzhu/configor)
- [stretchr/testify](https://github.com/stretchr/testify)
- [natefinch/lumberjack](https://https://github.com/natefinch/lumberjack)
- [prometheus/client_golang](https://github.com/prometheus/client_golang)
- [nacos-group/nacos-sdk-go](https://github.com/nacos-group/nacos-sdk-go) (可选)
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) (可选)

#### 命令行工具  

- [swag](https://github.com/swaggo/swag) (version>=1.8.4)
- [redoc-cli](https://github.com/Redocly/redoc) (version=latest)
- [golangci-lint](https://github.com/golangci/golangci-lint) (version>=1.47.0)  

## 如何使用？  

#### golang版本  
version >= 1.17.5

#### 启动服务  

在启动服务前，需要搭建必要依赖服务，如mysql数据库和redis缓存。为了帮助你快速的将服务运行起来，go-sail提供了基于docker-compose的基础服务启动脚本。
具体内容参考`ecosystem`目录下的相关内容。  
配置文件中的ip地址是随机样例，**实际值请修改成你自己的ip地址**（请**不要**使用`127.0.0.1`）。
- mysql服务  
```shell
cd ecosystem/docker-compose/mysql

docker-compose up -d
```  
命令执行后，将启动mysql服务，监听`33060`端口，账号/密码为：`root`/`root`。  

- redis服务  
> 命令执行前，请将`ecosystem/docker-compose/redis/docker-compose.yml`中的`192.168.224.114`全局替换为你自己的ip地址。
```shell
cd ecosystem/docker-compose/redis

docker-compose up -d
```  
命令执行后，将启动redis集群服务，集群以`cluster`模式运行，监听端口范围：`6379`~`6384`，认证密码为：`changeme`。  

- nacos服务（可选）  
```shell
cd ecosystem/docker-compose/nacos

docker-compose up -d
```  
命令执行后，将启动mysql服务和nacos服务，这里的mysql服务是独立的服务，旨在仅对nacos提供存储服务；nacos以`standalone`模式运行，账号/密码为：`nacos`/`nacos`。
浏览器访问`localhost:8848/nacos`，输入账号密码即可进入控制台。  
创建命名空间和配置文件：  
1.进入`命名空间`，点击`新增命名空间`，在`命名空间名`的输入框中输入`go-sail-user`，`描述`输入框中输入`go-sail user服务`，点击`确定`保存。  
2.进入`配置管理`>`配置列表`，在右侧选择`go-sail-user`命名空间，点击右侧的+号新增配置。在`Data ID`输入框中输入`go-sail-user.yml`，在`group`输入框中输入`go-sail`，
`配置格式`选择`YAML`，然后将`config-user.sample.yml`中的内容复制并粘贴到`配置内容`文本域中，点击发布。
> 记得将配置内容中的ip地址全局替换为你自己的ip地址。  

3.回到`命名空间`列表，记录下`go-sail-user`的命名空间id。  
4.设置环境变量并启动服务：  
```shell
export nacosAddrs=<nacos服务的地址> # 如：192.168.224.114:8848
export nacosNamespaceID=<go-sail-user的命名空间id>

go mod tidy

go run main.go user
```  
5.如果你不想使用nacos，也可以从本地配置文件启动。  
> 如果不从nacos读取配置启动服务，go-sail不会将服务注册到nacos中。

```shell
go mod tidy

go run main.go user -c ./config-user.sample.yml
```  

访问页面：  
- Swagger ui(debug=true)  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)  
- Redoc ui(debug=true)  
[http://localhost:8080/redoc/apidoc.html](http://localhost:8080/redoc/apidoc.html)  
- Prometheus metrics  
[http://localhost:1910/metrics](http://localhost:1910/metrics)  
- pprof(debug=true)  
[http://localhost:8080/debug/pprof](http://localhost:8080/debug/pprof)  

#### 命令脚手架  
- 生成openapi  
> 如果你的系统是Linux或MacOS，可直接使用make命令，更多指令请参考`Makefile`文件。  

```shell
make gen-swag-user
```  

- 构建镜像  
```shell
docker build --tag go-sail:v1.0.0 .
```  

- 本地构建  
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./go-sail
```

## 使用案例  
<img src="https://assets.pikaster-metaland.com/web/homepage/assets/img/logo.png" alt="Pikaster" width="600" />
<img src="https://www.wingoal-metaland.io/assets/images/text_a.png" alt="WinGoal" width="600" />

