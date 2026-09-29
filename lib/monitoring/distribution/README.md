# 分布式分位数统计

本目录提供跨服务实例的分布式数值分布统计能力，主要用于计算 HTTP 请求耗时等指标的 P85、P99、P99.9。实现采用 DDSketch 保存可合并的近似分布，同时精确保存样本数、总和、最小值和最大值；各节点只上报摘要，不上报原始样本。

## 统计原理

每条指标序列由 `namespace + service + metric + labels` 唯一确定，这些字段经过稳定排序和 SHA-256 计算得到 `MetricKey`。

每个节点在 `LocalWindow` 时间窗口内把观测值写入 DDSketch。窗口关闭后，节点生成一条 `SketchRecord`，其中包括：

- DDSketch 编码后的分布摘要；
- 精确的 `Count`、`Sum`、`Min` 和 `Max`；
- 节点 ID、指标信息、本地窗口边界和相对精度；
- 根据指标、节点和窗口生成的稳定 `RecordID`，用于幂等上报。

聚合节点按 `AggregateWindow` 查询所有完全落在目标窗口内的节点记录，并按 `MetricKey` 分组。DDSketch 支持直接合并，因此无需收集原始样本即可得到整个服务的近似分位值。聚合结果中的样本数、总和、最小值、最大值和平均值由节点精确摘要计算，只有分位值是近似值。

DDSketch 的 `RelativeAccuracy` 表示相对误差上限。例如配置为 `0.01` 时，分位值目标相对误差约为 1%。精度越高（数值越小），通常需要越多内存和存储空间。所有参与同一聚合的记录必须使用相同的相对精度。

时间窗口均采用左闭右开的形式 `[start, end)`，并按 Unix 时间边界对齐：

```text
业务观测
   │
   ▼
各节点 Registry / Recorder
   │  每个 LocalWindow 生成一条 DDSketch 记录
   ▼
有界上报队列 → ReportWorker（最多尝试 3 次）
   │
   ▼
Redis：节点记录、窗口索引
   │
   ▼
Redis 租约选出的唯一 Aggregator
   │  等待 AllowedLateness 后按 AggregateWindow 合并
   ▼
Redis：AggregatedResult + 聚合游标
   │
   ▼
Result / P85 / P99 / P999 查询
```

## 处理流程

1. `Registry.Observe` 或 `ObserveDuration` 将样本写入当前指标序列的本地 DDSketch。
2. `Registry` 每秒检查窗口；观测发生在窗口边界之后时，也会先轮转旧窗口，避免新样本落入旧窗口。
3. 已关闭窗口被编码为 `SketchRecord`，放入有界、非阻塞的上报队列。
4. `ReportWorker` 消费队列并写入 `Reporter`。单条记录失败时最多尝试 3 次，重试间隔依次为 1 秒、2 秒。
5. 启用聚合的实例通过 Redis 租约选主。只有持有租约的实例执行聚合，其他实例待命。
6. 当聚合窗口结束时间早于 `当前时间 - AllowedLateness` 时，Leader 合并窗口内的节点记录，计算配置的分位点并保存 `AggregatedResult`。
7. 成功聚合后推进持久化游标。Leader 重启或切换后会从游标继续补算遗漏窗口。
8. 上下文取消时，`Registry` 刷新当前非空窗口并关闭队列；Registry 和 ReportWorker 分别最多等待约 10 秒完成关闭阶段的入队和上报。

首次建立聚合游标时只处理最新的已关闭窗口，不会自动回溯更早的历史记录。游标建立后，停机期间遗漏的窗口才会被依次补算。

## 依赖组件

### DDSketch

通过 `github.com/DataDog/sketches-go/ddsketch` 实现近似分位数统计和跨节点摘要合并。当前传输结构版本为 `1`，算法标识为 `ddsketch-v1`。

### Redis

当前 `sail.Monitoring(...).Distribution(...)` 集成固定使用 `github.com/go-redis/redis/v8` 的 `UniversalClient`，可连接单机、哨兵或集群 Redis。Redis 用于：

- 幂等保存节点级 Sketch 记录；
- 建立本地窗口到记录 ID 的索引；
- 保存聚合结果和聚合游标；
- 通过 `SET NX`、TTL 和 Lua 脚本实现带所有权令牌的 Leader 租约。

Redis Key 前缀为 `go-sail:distribution`。节点记录与索引使用 `Reporter.Retention`，聚合结果使用 `Aggregator.ResultRetention`，聚合游标不会过期。

核心接口 `Reporter`、`Store` 和 `Leadership` 与具体后端解耦，可实现其他存储或租约后端；但 Go-Sail 的便捷接入目前仅组装 Redis 后端。

### 日志与生命周期

组件使用 `go.uber.org/zap` 记录上报、选主和聚合错误，并由调用方传入的 `context.Context` 控制退出。

## 配置

完整 YAML 示例见 [config.example.yaml](config.example.yaml)，同目录也提供 JSON 和 TOML 示例。

```yaml
distribution:
  enabled: true
  namespace: production
  service: user-center
  nodeID: user-center-01

  localWindow: 10s
  aggregateWindow: 1m
  allowedLateness: 20s
  relativeAccuracy: 0.01
  quantiles: [0.85, 0.99, 0.999]

  reporter:
    driver: redis
    queueSize: 1024
    retention: 24h

  aggregator:
    enabled: true
    leaderKey: go-sail:distribution:user-center:leader
    leaderTTL: 15s
    scanInterval: 5s
    resultRetention: 168h
```

配置包含顶层 `Conf`、上报配置 `ReporterConfig` 和聚合配置 `AggregatorConfig`。通过 `sail.Monitoring(...).Distribution(...)` 或 `NewRegistry(...)` 初始化时会先调用 `SetDefaults`，因此表中的默认值会写回传入的配置对象。

### 基础配置（`Conf`）

| 配置项 | Go 类型 | 必填 | 默认值 | 说明与约束 |
| --- | --- | --- | --- | --- |
| `enabled` | `bool` | 否 | `false` | 分布式统计总开关。为 `false` 时 `sail` 不创建 Registry、Reporter 或 Aggregator，并跳过其余字段校验。 |
| `namespace` | `string` | 启用时是 | 无 | 环境或租户隔离标识，参与生成 `MetricKey`。共享 Redis 的不同环境应使用不同值。 |
| `service` | `string` | 启用时是 | 无 | 产生指标的服务标识，参与生成 `MetricKey`、窗口索引和聚合游标。 |
| `nodeID` | `string` | 启用时是 | 无 | 当前服务实例的唯一标识，参与生成 `RecordID` 和统计 `NodeCount`。同一服务的不同实例不得复用。 |
| `localWindow` | `time.Duration` | 否 | `10s` | 单节点生成一条 Sketch 的窗口长度，必须大于 0。值越小，上报与 Redis 写入越频繁。 |
| `aggregateWindow` | `time.Duration` | 否 | `1m` | 跨节点结果窗口，必须大于 0，并且能被 `localWindow` 整除。 |
| `allowedLateness` | `time.Duration` | 否 | `20s` | 聚合窗口结束后继续等待迟到记录的时间。非正值会恢复为默认值，因此当前配置不能用 `0` 关闭等待。 |
| `relativeAccuracy` | `float64` | 否 | `0.01` | DDSketch 的相对误差上限，必须位于开区间 `(0, 1)`；同一服务所有节点必须一致。非正值会恢复为默认值。 |
| `quantiles` | `[]float64` | 否 | `[0.85, 0.99, 0.999]` | 聚合时预计算的分位点，每项必须位于 `(0, 1)`。初始化时按升序排序；当前实现不自动去重。只有这里配置的分位点才能通过结果的 `Quantile`/`P` 方法命中。 |
| `reporter` | `ReporterConfig` | 否 | 零值结构 | 节点记录的上报、内存队列及 Redis 保留策略，字段见下表。 |
| `aggregator` | `AggregatorConfig` | 否 | 零值结构 | Leader 选举、扫描和结果保留策略，字段见下表。 |

### 上报配置（`reporter`）

| 配置项 | Go 类型 | 必填 | 默认值 | 说明与约束 |
| --- | --- | --- | --- | --- |
| `reporter.driver` | `string` | 否 | 空字符串 | 预留的上报后端名称，例如 `redis`。当前 `sail.Monitoring(...).Distribution(...)` 固定组装 Redis Store，既不会根据该字段选择后端，也不会校验其值；建议仍显式填写 `redis` 以表达意图。 |
| `reporter.queueSize` | `int` | 否 | `1024` | Registry 到 ReportWorker 的内存队列容量。队列为非阻塞模式，满载时新关闭的窗口记录会被丢弃，并计入 `DroppedRecords`。非正值使用默认值。 |
| `reporter.retention` | `time.Duration` | 否 | `24h` | Redis 中节点级 Sketch、窗口索引和已聚合标记的 TTL。应覆盖聚合窗口、迟到等待和最长预期故障恢复时间。非正值使用默认值。 |

### 聚合配置（`aggregator`）

| 配置项 | Go 类型 | 必填 | 默认值 | 说明与约束 |
| --- | --- | --- | --- | --- |
| `aggregator.enabled` | `bool` | 否 | `false` | 当前实例是否参与 Leader 选举并运行聚合器。关闭它不影响当前实例采集和上报节点记录。部署中至少应有一个实例启用。 |
| `aggregator.leaderKey` | `string` | 否 | `go-sail:distribution:aggregator:leader` | Redis 租约键。参与同一服务聚合的候选节点必须相同；不同环境或服务必须使用不同 Key，建议包含 namespace 和 service。空字符串使用默认值。 |
| `aggregator.leaderTTL` | `time.Duration` | 否 | `15s` | Leader 租约有效期，必须严格大于 `scanInterval`。Leader 会在扫描周期内续租；非正值使用默认值。 |
| `aggregator.scanInterval` | `time.Duration` | 否 | `5s` | 尝试获取/续期租约以及扫描可聚合窗口的周期。它也决定窗口满足迟到等待后，结果最多还需等待多久才会生成。非正值使用默认值。 |
| `aggregator.resultRetention` | `time.Duration` | 否 | `168h`（7 天） | Redis 聚合结果的 TTL，不影响不会过期的聚合游标。非正值使用默认值。 |

### 配置格式与生效规则

- YAML 中 `time.Duration` 使用 `10s`、`1m`、`24h` 等字符串。
- JSON 和 TOML 中 `time.Duration` 使用纳秒整数，例如 `10s` 为 `10000000000`。
- 字段名区分大小写时应以结构标签为准，例如 `nodeID`、`localWindow`、`leaderTTL` 和 `resultRetention`。
- `SetDefaults` 会直接修改传入配置，并对 `quantiles` 原地排序。配置对象初始化后不应再被其他 goroutine 修改。
- `Validate` 会检查必填字段、窗口整除关系、分位点范围、相对精度范围以及 `leaderTTL > scanInterval`；当前不会校验 `reporter.driver`、保留期之间的关系或 `nodeID` 在集群中的实际唯一性。

## 使用方式

以下两种方式使用相同的统计模型和 Redis 数据结构。独立组件方式适合非 Go-Sail 项目、自定义指标或替换存储后端；Go-Sail 集成方式会自动接入框架的 HTTP 响应链路。

### 方式一：作为独立统计组件集成

直接使用时，应用负责创建并启动 `Registry`、`ReportWorker` 和 `Aggregator`：

- `Registry` 管理指标序列、本地窗口和上报队列；
- `ReportWorker` 消费 `Registry.ReportQueue()` 并调用 `Reporter`；
- `Aggregator` 通过 `Store` 和 `Leadership` 聚合已关闭窗口；
- `redisbackend.Store`、`redisbackend.Leadership` 是内置 Redis 实现；
- `exporter.Exporter`、`MultiExporter` 可用于应用自行发布已读取的终态结果，当前不会被 Aggregator 自动调用。

#### 初始化与启动

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

logger := zap.NewNop()
redisClient := redis.NewClient(&redis.Options{
    Addr: "127.0.0.1:6379",
})

conf := &distribution.Conf{
    Enabled:          true,
    Namespace:        "production",
    Service:          "user-center",
    NodeID:           "user-center-01",
    LocalWindow:      10 * time.Second,
    AggregateWindow:  time.Minute,
    AllowedLateness:  20 * time.Second,
    RelativeAccuracy: 0.01,
    Quantiles:        []float64{0.85, 0.99, 0.999},
    Reporter: distribution.ReporterConfig{
        Driver:    "redis",
        QueueSize: 1024,
        Retention: 24 * time.Hour,
    },
    Aggregator: distribution.AggregatorConfig{
        Enabled:         true,
        LeaderKey:       "go-sail:distribution:production:user-center:leader",
        LeaderTTL:       15 * time.Second,
        ScanInterval:    5 * time.Second,
        ResultRetention: 7 * 24 * time.Hour,
    },
}

// NewRegistry 会填充默认值并校验配置。
registry, err := distribution.NewRegistry(conf)
if err != nil {
    panic(err)
}

store := redisbackend.NewStore(
    redisClient,
    conf.LocalWindow,
    conf.Reporter.Retention,
    conf.Aggregator.ResultRetention,
)
leadership, err := redisbackend.NewLeadership(
    redisClient,
    conf.Aggregator.LeaderKey,
    conf.Aggregator.LeaderTTL,
)
if err != nil {
    panic(err)
}

worker := distribution.NewReportWorker(store, registry.ReportQueue(), logger)
aggregator := distribution.NewAggregator(conf, store, leadership, logger)

go registry.Start(ctx)
go worker.Run(ctx)
go aggregator.Run(ctx)
```

所有后台组件应共用同一个生命周期上下文。应用关闭时调用 `cancel()`，并为 Registry 刷新当前窗口、ReportWorker 排空队列预留足够时间。

#### 写入自定义指标

```go
labels := distribution.Labels{
    "operation": "create-user",
    "result":    "success",
}

// Observe 接受有限的非负数值，当前记录协议统一将数值解释为微秒。
registry.Observe("user.create.duration", 1250.5, labels)

// 耗时会自动转换为微秒，并保留小数部分。
startedAt := time.Now()
handleRequest()
registry.ObserveDuration("user.create.duration", time.Since(startedAt), labels)
```

`Observe` 和 `ObserveDuration` 不向调用方返回错误。无效样本计入 `Registry.ObserveErrors()`；窗口记录因队列满而丢失时计入 `Registry.DroppedRecords()`。

#### 查询聚合结果

底层组件通过 `MetricKey` 和精确窗口边界构造 `ResultID`，再从 Store 读取结果：

```go
windowEnd := time.Now().UTC().Truncate(conf.AggregateWindow)
windowStart := windowEnd.Add(-conf.AggregateWindow)

metricKey := distribution.BuildMetricKey(
    conf.Namespace,
    conf.Service,
    "user.create.duration",
    labels,
)
resultID := distribution.BuildResultID(metricKey, windowStart, windowEnd)

result, err := store.GetResult(ctx, resultID)
if err != nil {
    // 结果尚未生成、已过期或 Redis 访问失败。
    panic(err)
}
p99, exists := result.P99()
```

如需替换 Redis，可自行实现 `Reporter`、`Store` 和 `Leadership` 接口，并保持 `RecordID`、`ResultID` 的幂等语义。

### 方式二：在 Go-Sail 中以单例方式使用

Go-Sail 采用配置驱动方式管理分布式统计组件。启用 `monitor_conf` 后，框架在组件启动阶段自动创建全局 Monitor，组装 Registry、Redis Store、ReportWorker 和 Aggregator，并将 Registry 连接到 HTTP 响应链路。业务代码不需要手动调用 `sail.Monitoring(...).Distribution(...)`，后续通过 `sail.GetMonitor()` 获取单例即可。

#### 配置并启动

分布式统计依赖 Go-Sail 的 Redis 单机或集群组件，必须同时启用至少一种 Redis 配置。以下为 YAML 配置示例：

```yaml
redis_conf:
  enable: true
  endpoint:
    host: 127.0.0.1
    port: 6379
    username: ""
    password: ""
  database: 0
  ssl_enable: false

monitor_conf:
  enabled: true
  namespace: production
  service: user-center
  nodeID: user-center-01
  localWindow: 10s
  aggregateWindow: 1m
  allowedLateness: 20s
  relativeAccuracy: 0.01
  quantiles: [0.85, 0.99, 0.999]

  reporter:
    driver: redis
    queueSize: 1024
    retention: 24h

  aggregator:
    enabled: true
    leaderKey: go-sail:distribution:production:user-center:leader
    leaderTTL: 15s
    scanInterval: 5s
    resultRetention: 168h
```

配置加载到 `sail/config.Config` 后，按正常方式启动 Go-Sail 即可，不需要额外的 Monitor 初始化代码：

```go
sail.WakeupHttp("user-center", conf).
    Hook(registerRoutes, nil, nil).
    Launch()
```

框架启动和关闭行为如下：

- `componentsStartup` 先初始化 Logger 和 Redis，再根据 `MonitorDistributorConf.Enabled` 初始化 Monitor。
- 同时配置 Redis 单机和集群时，`sail.GetRedis()` 优先使用单机实例。
- Monitor 已启用但 Redis 未启用时，框架取消 Monitor 初始化并输出提示；此时 `sail.GetMonitor()` 返回 `nil`。
- Monitor 在进程内只初始化一次，不会重复启动 Registry、ReportWorker 或 Aggregator。
- `componentsShutdown` 会先取消 Monitor 的运行上下文，再关闭 Redis 等基础组件。

每次通过 `sail.Response(...).Send()` 或 `SendWithCode(...)` 返回 HTTP 响应时，系统自动记录：

- 指标名：`http.server.duration`；
- 单位：微秒，可保留小数；
- 标签：`method`、Gin 完整路由 `route`、状态码类别 `status`（如 `2xx`）。

#### 通过 `GetMonitor` 查询聚合结果

完成初始化后，通过 `sail.GetMonitor()` 获取全局实例即可在任意业务位置查询，无需传递或重新创建 Monitor：

```go
sailMonitor := sail.GetMonitor()
if sailMonitor == nil {
    return errors.New("monitoring is not initialized")
}

appConf := config.Get()
if appConf == nil {
    return errors.New("go-sail config is not initialized")
}
monitorConf := appConf.MonitorDistributorConf

windowEnd := time.Now().UTC().Truncate(monitorConf.AggregateWindow)
windowStart := windowEnd.Add(-monitorConf.AggregateWindow)
labels := distribution.Labels{
    "method": "GET",
    "route":  "/users/:id",
    "status": "2xx",
}

result, err := sailMonitor.Result(
    "http.server.duration",
    labels,
    windowStart,
    windowEnd,
)
if err != nil {
    // 结果尚未生成、已过期或 Redis 访问失败。
    return err
}

p99, exists, err := sailMonitor.P99(
    "http.server.duration",
    labels,
    windowStart,
    windowEnd,
)
if err != nil {
    return err
}
if exists {
    logger.Info("distribution result",
        zap.Float64("p99", p99),
        zap.Uint64("samples", result.SampleCount),
    )
}
```

`Result` 返回的主要字段包括 `NodeCount`、`RecordCount`、`SampleCount`、`Sum`、`Min`、`Max`、`Avg` 和 `Quantiles`。`P85`、`P99`、`P990`、`P999` 是便捷方法；其中 `P990` 表示 P99.0，与 P99 等价。也可以使用 `P(percentile, ...)` 查询已配置的分位点，例如 `P(999, ...)` 表示 P99.9。

结果通常要到 `windowEnd + allowedLateness` 之后的下一次扫描才可查询。例如窗口在 `12:01:00` 结束、迟到等待为 20 秒、扫描间隔为 5 秒，则结果一般不会早于 `12:01:20` 生成。

## 注意事项

1. **结果是窗口终态快照。** 聚合器超过迟到等待期后只处理一次窗口并推进游标；更晚到达的记录不会自动改写已经终态化的结果。`AllowedLateness` 应覆盖正常网络延迟、重试时间和短暂 Redis 抖动。
2. **上报队列会主动丢弃。** 正常运行期间入队是非阻塞的；队列满时记录会被丢弃，可通过 `Registry.DroppedRecords()` 观测。应结合峰值指标序列数调整 `queueSize`。
3. **观测错误不会返回。** 非有限值、负数或记录器错误由 `Registry.ObserveErrors()` 计数。调用方应将这两个内部计数接入现有监控。
4. **控制标签基数。** 每种标签组合都会创建独立 Recorder。不要使用用户 ID、请求 ID、原始 URL 等高基数字段；HTTP 路由应使用模板路径。
5. **标签内容必须一致。** 查询时标签的键和值必须与写入时完全相同。Registry 会复制标签，但调用方仍不应在提交后并发修改原映射。
6. **统一配置。** 同一 `namespace/service` 的实例应使用相同的窗口、相对精度和分位点配置，并使用唯一 `nodeID`。精度不一致的记录会被聚合器拒绝。
7. **Leader Key 必须隔离。** 默认 Leader Key 是全局固定值。多环境、多服务共享 Redis 时必须显式加入环境和服务名，否则不同服务会竞争同一把租约。
8. **至少部署一个聚合候选节点。** 可以在所有实例启用聚合，由租约保证单 Leader；也可以只在专用实例启用。所有候选节点必须连接同一 Redis。
9. **合理设置保留期。** `reporter.retention` 至少应大于聚合窗口、迟到等待和预期最长故障恢复时间之和，否则补算时原始记录可能已经过期。
10. **Redis 可用性影响完整性。** 上报只做有限次数重试，没有磁盘缓冲。Redis 长时间不可用会造成记录缺失；对完整性要求高时应实现持久化 Reporter 或外部队列。
11. **查询边界必须精确对齐。** `ResultID` 包含纳秒级窗口边界；任一边界不同都会查询到另一个 ID。建议统一用 UTC，并用 `Truncate(AggregateWindow)` 计算边界。
12. **当前导出器不会被自动调用。** `Exporter` 只是扩展接口，聚合器目前仅将结果写入 Store；如需推送到 Prometheus、消息队列或其他系统，需要应用自行读取结果并调用导出器。
13. **关闭要预留时间。** 取消上下文后组件会尽力刷新和排空，但不是无限等待。进程应提供足够的优雅退出时间。
