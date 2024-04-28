package logger

import (
	"fmt"
	"log"
	"os"
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
	module   string
	instance *zap.Logger
}

func (sc *svcHolders) load(module string) svcStd {
	for _, svc := range sc.instances {
		if svc.module == module {
			return svc
		}
	}

	//panic(fmt.Sprintf("GetLogger failure:logger instance [%s] not initialize", module))

	//使用默认module
	return sc.instances[0]
}

func (sc *svcHolders) store(module string, lg *zap.Logger) {
	var exist bool
	for k, svc := range sc.instances {
		//已存在，则更新
		if svc.module == module {
			sc.instances[k].instance = lg
			exist = true
			break
		}
	}

	//不存在，则新增
	if !exist {
		sc.instances = append(sc.instances, svcStd{module: module, instance: lg})
	}
}

var (
	gDefaultModeName  = "<defaultModeName>"
	gLoggerSvcHolders *svcHolders
	gWriterSyncers    = make([]zapcore.WriteSyncer, 0, 2)
)

// GetLogger 获取日志服务实例
//
// 要使用对应module，请初始化时指定 Conf.Modules 配置
//
// 如果module未被初始化，那么将使用默认module
func GetLogger(module ...string) *zap.Logger {
	if len(module) < 1 {
		return gLoggerSvcHolders.load(gDefaultModeName).instance
	}

	return gLoggerSvcHolders.load(module[0]).instance
}

// Init 初始化
//
// InitLoggerZap 方法的语法糖
func Init(cfg Conf, appName string, syncers ...zapcore.WriteSyncer) {
	InitLoggerZap(cfg, appName, syncers...)
}

// InitLoggerZap 初始化zap日志服务
//
// 会加入默认的一个模块空间，当不传参调用 GetLogger 时，
// 就是使用默认的模块空间
//
// 当启用exporter时，Exporter.Provider 字段值将作为exporter方案，
// 为空则表示不启用exporter。
// 目前Provider字段支持:
//
// 'redis'、'redis-cluster'、'nats'、'kafka'
func InitLoggerZap(cfg Conf, appName string, syncers ...zapcore.WriteSyncer) {
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

		//输出到文件
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(fileWriter)), atomicLevel,
		)

		cores = append(cores, fileCore)

		//输出到终端(如果配置启用)
		if cfg.ConsoleOutput {
			//consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
			consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

			consoleDebugging := zapcore.Lock(os.Stdout)
			consoleErrors := zapcore.Lock(os.Stderr)

			consoleCore := zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, consoleErrors, zapcore.ErrorLevel),
				zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
			)

			cores = append(cores, consoleCore)
		}

		//logstash订阅的key只定义一个，与module无关
		writer := exporterProvider(cfg)
		if writer != nil {
			coreWithWriter := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(writer)), atomicLevel,
			)
			cores = append(cores, coreWithWriter)
		}

		//设置自定义exporter
		if len(syncers) > 0 {
			gWriterSyncers = append(gWriterSyncers, syncers...)
		}

		//读取外部配置的syncer并加入到cores中
		for _, syncer := range gWriterSyncers {
			coreWithWriter := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(syncer)), atomicLevel,
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
		log.Println("[logger] using (redis) exporter")
		return writer
	case "redis-cluster":
		redisWriter := &redisClusterWriterStd{
			cli:     redis.NewCluster(cfg.Exporter.Redis.ClusterConnConf),
			listKey: cfg.Exporter.Redis.ListKey,
		}

		writer = redisWriter
		log.Println("[logger] using (redis-cluster) exporter")
		return writer
	case "nats":
		natsWriter := &natsWriterStd{
			cli:        nats.New(cfg.Exporter.Nats.ConnConf),
			subjectKey: cfg.Exporter.Nats.Subject,
		}

		writer = natsWriter
		log.Println("[logger] using (nats) exporter")
		return writer
	case "kafka":
		kafkaWriter := &kafkaWriterStd{
			writer: kafka.NewWriter(cfg.Exporter.Kafka.ConnConf, cfg.Exporter.Kafka.Topic),
			topic:  cfg.Exporter.Kafka.Topic,
		}

		writer = kafkaWriter
		log.Println("[logger] using (kafka) exporter")
		return writer
	default:
		log.Println("[logger] exporter not set,ignore emit exporter")
		return writer
	}
}
