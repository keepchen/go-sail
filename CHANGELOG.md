# Change Log  

## [v3.0.6-rc3] – 2025-09-04

✨ 新增 valkey 组件、支持自定义 redis 客户端、修复多处空指针问题，升级 jwt/gopsutil 等依赖

<details open>
<summary>中文</summary>

#### 🚀 新功能
- utils: 简化并优化 redis 锁代码，支持传入自定义 redis 客户端   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- schedule: 支持设定自定义 redis 客户端   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- http: 调整响应器时间戳字段赋值位置，分页器 JSON tag 修正   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- other: 更新 README ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- schedule: 新增 crontab 表达式；utils: HTTP 请求不再检测响应状态码 ([`9be8d24`](https://github.com/keepchen/go-sail/commit/9be8d24))
- lib: db 组件新增 `NowFunc` 配置 ([`d7b1f79`](https://github.com/keepchen/go-sail/commit/d7b1f79))
- middleware: 新增获取客户端真实 IP 方法 ([`edb4b3a`](https://github.com/keepchen/go-sail/commit/edb4b3a))
- lib: 新增 valkey 组件 ([`c9a53b7`](https://github.com/keepchen/go-sail/commit/c9a53b7))
- lib: nacos 组件新增服务订阅方法，并替换旧 utils 方法调用 ([`8f8e793`](https://github.com/keepchen/go-sail/commit/8f8e793))
- lib: nacos新增`NewConfigClient`和`NewNamingClient`方法 ([`30b6307b`](https://github.com/keepchen/go-sail/commit/30b6307b))
- sail: 新增config配置文件读取 ([`c70b1c7e`](https://github.com/keepchen/go-sail/commit/c70b1c7e))
- sail: jwt新增加解密方法 ([`737b694c`](https://github.com/keepchen/go-sail/commit/737b694c))
- sail: 新增`RedisLocker`方法调用 ([`31c55834`](https://github.com/keepchen/go-sail/commit/31c55834))
- sail: 新增setter统一管理redis锁和schedule的redis实例 ([`b20009a1`](https://github.com/keepchen/go-sail/commit/b20009a1))

#### 🐛 修复
- schedule: 修复 `Call` 和 `MustCall` 空指针问题 ([`ebd4ea9`](https://github.com/keepchen/go-sail/commit/ebd4ea9))
- api: 修复 `mergeBody` 对 `(*T)(nil)` 的处理问题 ([`817b93f`](https://github.com/keepchen/go-sail/commit/817b93f))
- api: 修复 `SendWithCode` 中 `funcBeforeWrite` 空指针问题 ([`f41fa8f`](https://github.com/keepchen/go-sail/commit/f41fa8f))

#### 🔧 变更 / 优化
- middleware: gopsutil 升级到 v4 ([`6584811`](https://github.com/keepchen/go-sail/commit/6584811))
- lib: jwt 修改错误文案 ([`3b396e4`](https://github.com/keepchen/go-sail/commit/3b396e4))
- sail: jwt `ValidToken` 返回参数调整 ([`2ab099e`](https://github.com/keepchen/go-sail/commit/2ab099e))
- lib: jwt 升级到 v5 ([`5a94765`](https://github.com/keepchen/go-sail/commit/5a94765))
- utils: redis 锁调整 `TryLockWithContext` 方法 ([`c578ab4`](https://github.com/keepchen/go-sail/commit/c578ab4))
- other: 框架版本号修改为`3.0.6`  
- http: api响应器性能优化 ([`23934799`](https://github.com/keepchen/go-sail/commit/23934799))  
- other: 新增benchmark ([`23934799`](https://github.com/keepchen/go-sail/commit/23934799))

#### 📦 依赖升级
- github.com/golang-jwt/jwt/v5 → 5.2.2 → 5.3.0
- github.com/shirou/gopsutil/v4 → 4.25.3 → 4.25.7
- 其他依赖升级：swag, etcd, gorm, mysql, sqlite, postgres, nats, gin, nacos, valkey, kafka, x/net 等

#### 📖 文档 & 🧪 测试
- 更新 README / README_EN.md / examples
- 新增测试用例 & codecov 配置
- 持续完善测试用例 & CI/CD workflow 调整（多个提交）

</details>

---

<details open>
<summary>English</summary>

#### 🚀 Features
- utils: Simplified and optimized redis lock code, support custom redis client   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- schedule: Support custom redis client   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- http: Adjusted timestamp field, fixed paginator JSON tag   ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- other: Updated README ([`0166990`](https://github.com/keepchen/go-sail/commit/0166990))
- schedule: Added crontab expression; utils: HTTP requests no longer check status code ([`9be8d24`](https://github.com/keepchen/go-sail/commit/9be8d24))
- lib: Added `NowFunc` in db component ([`d7b1f79`](https://github.com/keepchen/go-sail/commit/d7b1f79))
- middleware: Added real client IP method ([`edb4b3a`](https://github.com/keepchen/go-sail/commit/edb4b3a))
- lib: Added valkey component ([`c9a53b7`](https://github.com/keepchen/go-sail/commit/c9a53b7))
- lib: nacos: Added service subscription, replaced old utils calls ([`8f8e793`](https://github.com/keepchen/go-sail/commit/8f8e793))
- lib: nacos: Added `NewConfigClient` and `NewNamingClient` methods ([`30b6307b`](https://github.com/keepchen/go-sail/commit/30b6307b))
- sail: Added `Config` to read configuration ([`c70b1c7e`](https://github.com/keepchen/go-sail/commit/c70b1c7e))
- sail: jwt: Added `Encrypt` and `Decrypt` methods ([`737b694c`](https://github.com/keepchen/go-sail/commit/737b694c))
- sail: Added `RedisLocker` method ([`31c55834`](https://github.com/keepchen/go-sail/commit/31c55834))
- sail: Added `setter` to manage redis client for redis locker and schedule ([`b20009a1`](https://github.com/keepchen/go-sail/commit/b20009a1))

#### 🐛 Fixes
- schedule: Fixed `Call` and `MustCall` nil pointer issue ([`ebd4ea9`](https://github.com/keepchen/go-sail/commit/ebd4ea9))
- api: Fixed `mergeBody` handling for `(*T)(nil)` ([`817b93f`](https://github.com/keepchen/go-sail/commit/817b93f))
- api: Fixed `SendWithCode` nil pointer when calling `funcBeforeWrite` ([`f41fa8f`](https://github.com/keepchen/go-sail/commit/f41fa8f))

#### 🔧 Changes / Improvements
- middleware: Upgraded gopsutil to v4 ([`6584811`](https://github.com/keepchen/go-sail/commit/6584811))
- lib: Modified jwt error messages ([`3b396e4`](https://github.com/keepchen/go-sail/commit/3b396e4))
- sail: Adjusted jwt `ValidToken` return parameters ([`2ab099e`](https://github.com/keepchen/go-sail/commit/2ab099e))
- lib: Upgraded jwt to v5 ([`5a94765`](https://github.com/keepchen/go-sail/commit/5a94765))
- utils: Redis lock adjusted `TryLockWithContext` ([`c578ab4`](https://github.com/keepchen/go-sail/commit/c578ab4))
- other: change Framework version to `3.0.6`  
- http: Api performance optimize ([`23934799`](https://github.com/keepchen/go-sail/commit/23934799))
- other: Add benchmark ([`23934799`](https://github.com/keepchen/go-sail/commit/23934799))

#### 📦 Dependencies
- github.com/golang-jwt/jwt/v5 → 5.2.2 → 5.3.0
- github.com/shirou/gopsutil/v4 → 4.25.3 → 4.25.7
- Other deps: swag, etcd, gorm, mysql, sqlite, postgres, nats, gin, nacos, valkey, kafka, x/net

#### 📖 Docs & 🧪 Tests
- Updated README / README_EN.md / examples
- Added test cases & codecov config
- Continuous test improvements & CI/CD workflow adjustments

</details>

## v3.0.5_rc  
### 常量
- 新增测试用例 ([`def9f3d2`](https://github.com/keepchen/go-sail/commit/def9f3d2))  ([`46df549b`](https://github.com/keepchen/go-sail/commit/46df549b))

### 类库
- 数据库日志配置调整 ([`a8a4cec4`](https://github.com/keepchen/go-sail/commit/a8a4cec4))
- 数据库日志打印规则调整 ([`a8a4cec4`](https://github.com/keepchen/go-sail/commit/a8a4cec4))
- Jwt新增`MergeStandardClaims`方法 ([`f48f0f3d`](https://github.com/keepchen/go-sail/commit/f48f0f3d))
- 数据库新增Gorm配置项 ([`27002b8b`](https://github.com/keepchen/go-sail/commit/27002b8b))
- 新增Notification库，支持lark、dingtalk、和slack ([`29c29ac8`](https://github.com/keepchen/go-sail/commit/29c29ac8)) ([`773b1277`](https://github.com/keepchen/go-sail/commit/773b1277))
- Nacos的`InitClient`函数新增客户端入参 ([`d1efc9f0`](https://github.com/keepchen/go-sail/commit/d1efc9f0))

### 路由中间件
- 跨域中间件不再特别返回204状态码 ([`5487e0b1`](https://github.com/keepchen/go-sail/commit/5487e0b1))
- 新增限流器中间件 ([`41fe9b7c`](https://github.com/keepchen/go-sail/commit/41fe9b7c))
- Prometheus新增系统指标采样 ([`138a6a20`](https://github.com/keepchen/go-sail/commit/138a6a20))
- 限流器新增Redis支持 ([`634d0cbf`](https://github.com/keepchen/go-sail/commit/634d0cbf))
- 限流器本地方案改用sync.Map提升性能 ([`634d0cbf`](https://github.com/keepchen/go-sail/commit/634d0cbf))
- 新增限流器测试用例 ([`634d0cbf`](https://github.com/keepchen/go-sail/commit/634d0cbf))
- `WithCorsOnlyOptions`方法返回200状态码  ([`fd55ae31`](https://github.com/keepchen/go-sail/commit/fd55ae31))
- Logtrace新增最长requestId限制 ([`fd55ae31`](https://github.com/keepchen/go-sail/commit/fd55ae31))

### 计划任务
- 新增手动调用语法糖 ([`6234e332`](https://github.com/keepchen/go-sail/commit/6234e332))

### 工具类
- 新增Number相关方法 ([`b52fc1af`](https://github.com/keepchen/go-sail/commit/b52fc1af))
- Redis分布式锁的值修改为持有者信息 ([`ddcc80d8`](https://github.com/keepchen/go-sail/commit/ddcc80d8))
- 重写随机浮点数方法 ([`0c0e2f49`](https://github.com/keepchen/go-sail/commit/0c0e2f49))
- Redis锁代码优化 ([`0c0e2f49`](https://github.com/keepchen/go-sail/commit/0c0e2f49))
- 新增字符串/字节数组转换函数 ([`b52c2169`](https://github.com/keepchen/go-sail/commit/b52c2169))
- 新增`SendRequest`方法 ([`b1c2766a`](https://github.com/keepchen/go-sail/commit/b1c2766a))
- 新增Gzip压缩/解压方法 ([`b1c2766a`](https://github.com/keepchen/go-sail/commit/b1c2766a))
- 新增Domain和Cert工具函数 ([`29c29ac8`](https://github.com/keepchen/go-sail/commit/29c29ac8))
- **对工具类进行分组改造** ([`f6bf3181`](https://github.com/keepchen/go-sail/commit/f6bf3181))
- **原工具类方法标记为废弃** ([`f6bf3181`](https://github.com/keepchen/go-sail/commit/f6bf3181))
- RSA新增格式化函数兼容多种格式的公私钥 ([`3f31b32e`](https://github.com/keepchen/go-sail/commit/3f31b32e))

### 响应器
- dto.Base修改swagger注释 ([`29ea7a4d`](https://github.com/keepchen/go-sail/commit/29ea7a4d))
- Api设置新增`FuncBeforeWrite`函数  ([`29ea7a4d`](https://github.com/keepchen/go-sail/commit/29ea7a4d))
- Api写入响应前调用`FuncBeforeWrite`函数 ([`f2ed64c4`](https://github.com/keepchen/go-sail/commit/f2ed64c4))
- 修改分页实体tag ([`138a6a20`](https://github.com/keepchen/go-sail/commit/138a6a20))
- 设置项新增更多内置错误码覆盖选项 ([`e15dd1d3`](https://github.com/keepchen/go-sail/commit/e15dd1d3))

### ORM
- 代码优化  ([`b52fc1af`](https://github.com/keepchen/go-sail/commit/b52fc1af))
- 新增Hook时间支持 ([`b52fc1af`](https://github.com/keepchen/go-sail/commit/b52fc1af))
- 新增`NewSvcImplSilent`方法 ([`6234e332`](https://github.com/keepchen/go-sail/commit/6234e332))
- `BeforeSave`加入空值检查 ([`b1a7f0f6`](https://github.com/keepchen/go-sail/commit/b1a7f0f6))

### 框架
- sail关键字新增marshal日志支持 ([`b52fc1af`](https://github.com/keepchen/go-sail/commit/b52fc1af))
- 新增`GetRedisUniversal`方法 ([`b1c2766a`](https://github.com/keepchen/go-sail/commit/b1c2766a))
- `GetRedis`方法变更为获取通用实例 ([`b1c2766a`](https://github.com/keepchen/go-sail/commit/b1c2766a))
- 新增Logtrace相关方法 ([`cd8c71fc`](https://github.com/keepchen/go-sail/commit/cd8c71fc))
- 控制台打印信息新增仓库地址 ([`3e9daf3c`](https://github.com/keepchen/go-sail/commit/3e9daf3c))
- 新增Redis、Nats、Etcd、Kafka组件新增实例方法 ([`3e9daf3c`](https://github.com/keepchen/go-sail/commit/3e9daf3c))
- 新增Jwt相关语法糖 ([`fd55ae31`](https://github.com/keepchen/go-sail/commit/fd55ae31))
- 框架版本号更新到3.0.5

### 其他
- **将Go最低版本要求提升到`1.20`** ([`7772e680`](https://github.com/keepchen/go-sail/commit/7772e680))
- 框架终端打印方法代码优化 ([`66fcd085`](https://github.com/keepchen/go-sail/commit/66fcd085))
- 更新README文档
- 更新Examples调用示例
- Swagger文档路由加入空配置判断 ([`f48f0f3d`](https://github.com/keepchen/go-sail/commit/f48f0f3d))

## v3.0.4  
### Config
- jwt配置改为指针类型

### 常量
- 新增错误码注入方法`RegisterCodeSingle`和`RegisterCodeTable`
- 原错误码注入方法`RegisterCode`标记为弃用

### 类库
- nacos新增获取配置方法`GetConfig`
- nacos新增配置监听方法`ListenConfigWithCallback`
- nacos配置监听方法新增是否打印原字符参数
- nacos组件库日志等级调整为warn
- redis单实例配置tag修正
- logger配置注释修正
- logger新增终端输出支持
- logger初始化函数新增syncers可选参数以支持自定义导出器
- jwt验证签名不再从私钥解析公钥而是直接使用公钥
- jwt新增`MustLoad`方法，原`Load`方法逻辑变更为公私钥二者存在其一即可
- jwt中MapClaims的`Valid`方法继承`jwtLib.StandardClaims的Valid`
- 部分组件`New`方法出现错误不再panic而是返回错误
- etcd新增服务注册与发现方法
- redis去除无用配置代码

### 路由中间件
- Websocket新增中间件支持
- 跨域中间件加入请求方法判断
- Prometheus中间件加入重入检测

### 计划任务
- 新增语法糖`EveryFifteenSeconds`,`EveryFifteenMinutes`
- 任务名重复时将panic
- 更新代码注释

### 工具类
- md5修改方法名
- redislock新增`XXWithContext`语法糖
- redislock代码优化
- **[Fix]** 重写随机浮点数方法
- 新增heap操作
- 时间工具新增语法糖
- 新增`FromCharCode`和`CharCodeAt`方法

### 响应器
- **[Fix]** 时区对象空指针修复
- dto.Base中的code类型变更为int
- dto.Base新增测试用例
- 新增`DefaultSetupOption`方法
- 新增调用方法并标记部分方法为弃用状态
- 空data字段处理逻辑

### ORM
- 此模块为新增模块

### 框架
- **[Fix]** 启动错误修复(空指针检测)
- 新增组件初始化成功提示
- 服务终止后按配置依次关闭组件
- Prometheus服务改为支持信号监测优雅退出
- http服务设置默认监听地址为':8080'
- 启动函数中的beforeFunc和afterFunc变更为异步执行
- 启动成功的终端信息打印新增`swagger ui`地址
- 框架版本更新到3.0.4

### 其他
- 更新README文档
- 更新examples调用示例
- 修改注释避免与swag解析冲突
- `.github`目录新增issue模板
-  新增`orm`模块

## v3.0.3  
1.Config
- 新增Set方法
- 新增解析配置到目标结构体方法`ParseConfigFromBytesToDst`

2.路由中间件
- `RequestEntry`中间件更名为`LogTrace`
- 跨域中间件新增`WithCorsOnlyOptions`

3.计划任务模块
- 新增`RunAfter`,`FirstDayOfWeek`,`LastDayOfWeek`方法

4.框架
- 错误恢复时打印调用堆栈
- 新增Websocket支持
- 框架版本号更新为3.0.3

5.其他优化
- 更新README
- 更新examples调用示例
- 重建.gitignore缓存
- 更换彩色Logo
- 代码优化  

## v3.0.2  
1.工具类
- redislock新增RedisTryLock方法
- redislock中RedisLock方法变更为阻塞式
- **[Fix]** redislock自动续期管理bug修复

2.响应器
- 错误码新增多语言支持

3.计划任务模块
- 新增状态查询支持
- 组件方法代码接口化调整

4.框架
- 框架启动方法优化
- 启动gin引擎时默认使用requestEntry中间件
- 组件方法代码接口化调整
- 框架版本号更新为3.0.2

5.其他优化
- 错误码新增多语言支持
- **[Fix]** sync.Once使用错误
- 代码优化  

## v3.0.1  
1.utils工具类新增方法
- sm4加解密
- md5摘要计算
- 软件版本打印
- 中国大陆身份证验证

2.lib组件库新增组件
- 新增etcd连接
- 新增kafka连接
- logger导出器新增kafka支持
- logger组件GetLogger方法加入modules参数支持
- 本地cache新增list链表操作支持

3.新增计划任务模块

4.路由中间件
- 新增浏览器客户端语言解析
- 请求入口中间件上下文注入新增spanId

5.框架
- 新增组件获取函数
- 启动函数新增before和after自定义函数
- 更改框架版本号

6.其他优化
- 代码注释统一
- 框架日志打印统一
- 更新readme
- 更新examples

## v3.0.0  
- Complete the framework transformation
- Optimize toolkit functions  

## v2.0.3  
- Optimize toolkit functions
- Fix typo  