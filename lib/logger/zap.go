package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/keepchen/go-sail/v3/lib/kafka"

	"github.com/keepchen/go-sail/v3/lib/nats"

	"github.com/keepchen/go-sail/v3/lib/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type svcHolders struct {
	instances []svcStd
}

type svcStd struct {
	modeName string
	instance *zap.Logger
}

func (sc *svcHolders) load(modeName string) svcStd {
	for _, svc := range sc.instances {
		if svc.modeName == modeName {
			return svc
		}
	}

	//panic(fmt.Sprintf("GetLogger failure:logger instance [%s] not initialize", modeName))

	//使用默认modeName
	return sc.instances[0]
}

func (sc *svcHolders) store(modeName string, lg *zap.Logger) {
	var exist bool
	for k, svc := range sc.instances {
		//已存在，则更新
		if svc.modeName == modeName {
			sc.instances[k].instance = lg
			exist = true
			break
		}
	}

	//不存在，则新增
	if !exist {
		sc.instances = append(sc.instances, svcStd{modeName: modeName, instance: lg})
	}
}

var (
	gDefaultModeName  = "<defaultModeName>"
	gLoggerSvcHolders *svcHolders
	gWriterSyncers    = make([]zapcore.WriteSyncer, 0, 1)
)

// GetLogger 获取日志服务实例
func GetLogger(modeName ...string) *zap.Logger {
	if len(modeName) < 1 {
		return gLoggerSvcHolders.load(gDefaultModeName).instance
	}

	return gLoggerSvcHolders.load(modeName[0]).instance
}

// SetExporters 设置导出器
//
// 设置自定义的导出器
func SetExporters(syncers []zapcore.WriteSyncer) {
	if len(syncers) > 0 {
		gWriterSyncers = append(gWriterSyncers, syncers...)
	}
}

// Init 初始化
//
// InitLoggerZap 方法的语法糖
func Init(cfg Conf, appName string) {
	InitLoggerZap(cfg, appName)
}

// InitLoggerZap 初始化zap日志服务
//
// 会加入默认的一个模块空间，当不传参调用GetLogger()时，
// 就是使用默认的模块空间
//
// 当启用elk时，logger根据provider配置使用redis队列或nats publish等作为媒介，需要在logstash侧配置对应的pipeline
// 队列的key取决于日志文件名和appName的组合，如：
// 日志文件名=logs/app.log，appName=app
// 则，队列名称为=> app:logs/app.log
func InitLoggerZap(cfg Conf, appName string) {
	//注入默认的空间模块
	cfg.Modules = append(cfg.Modules, gDefaultModeName)
	sc := &svcHolders{}

	//定义全局日志组件配置
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(cfg.Level) {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "dpanic":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "panic":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "fatal":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	for _, mn := range cfg.Modules {
		var (
			filename string
			cores    []zapcore.Core
		)
		if len(cfg.Filename) == 0 {
			cfg.Filename = "logs/running.log"
		}
		if mn != gDefaultModeName {
			filename = strings.Replace(cfg.Filename, ".log", fmt.Sprintf("_%s.log", mn), 1)
		} else {
			filename = cfg.Filename
		}

		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  true,
			Compress:   cfg.Compress,
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriter), atomicLevel,
		)

		cores = append(cores, fileCore)

		//logstash订阅的key只定义一个，与modeName无关
		writer := exporterProvider(cfg)
		if writer != nil {
			coreWithWriter := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), atomicLevel,
			)
			cores = append(cores, coreWithWriter)
		}

		//读取外部配置的syncer并加入到cores中
		for _, syncer := range gWriterSyncers {
			coreWithWriter := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(syncer), atomicLevel,
			)
			cores = append(cores, coreWithWriter)
		}

		zapCore := zapcore.NewTee(cores...)

		loggerWithFields := zap.New(zapCore, zap.AddCaller()).With(zap.String("serviceName", fmt.Sprintf("%s:%s", appName, mn)))

		sc.store(mn, loggerWithFields)
	}

	gLoggerSvcHolders = sc

	defer func() {
		for _, mn := range cfg.Modules {
			_ = gLoggerSvcHolders.load(mn).instance.Sync()
		}
	}()
}

// 设置导出器
func exporterProvider(cfg Conf) zapcore.WriteSyncer {
	var writer zapcore.WriteSyncer

	switch strings.ToLower(cfg.Exporter.Provider) {
	case "redis":
		redisWriter := &redisWriterStd{
			cli:     redis.New(cfg.Exporter.Redis.ConnConf),
			listKey: cfg.Exporter.Redis.ListKey,
		}

		writer = redisWriter
		log.Println("[logger] using (redis) writer")
		return writer
	case "redis-cluster":
		redisWriter := &redisClusterWriterStd{
			cli:     redis.NewCluster(cfg.Exporter.Redis.ClusterConnConf),
			listKey: cfg.Exporter.Redis.ListKey,
		}

		writer = redisWriter
		log.Println("[logger] using (redis-cluster) writer")
		return writer
	case "nats":
		natsWriter := &natsWriterStd{
			cli:        nats.New(cfg.Exporter.Nats.ConnConf),
			subjectKey: cfg.Exporter.Nats.Subject,
		}

		writer = natsWriter
		log.Println("[logger] using (nats) writer")
		return writer
	case "kafka":
		kafkaWriter := &kafkaWriterStd{
			writer: kafka.NewWriter(cfg.Exporter.Kafka.ConnConf, cfg.Exporter.Kafka.Topic),
			topic:  cfg.Exporter.Kafka.Topic,
		}

		writer = kafkaWriter
		log.Println("[logger] using (kafka) writer")
		return writer
	default:
		log.Println("[logger] writer not set,ignore emit exporter")
		return writer
	}
}
