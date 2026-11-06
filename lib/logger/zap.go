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

	//Use the default module
	return sc.instances[0]
}

func (sc *svcHolders) store(module string, lg *zap.Logger) {
	var exist bool
	for k, svc := range sc.instances {
		//Update if it exists
		if svc.module == module {
			sc.instances[k].instance = lg
			exist = true
			break
		}
	}

	//Add if it does not exist
	if !exist {
		sc.instances = append(sc.instances, svcStd{module: module, instance: lg})
	}
}

var (
	gDefaultModeName  = "<defaultModeName>"
	gLoggerSvcHolders *svcHolders
	gWriterSyncers    = make([]zapcore.WriteSyncer, 0, 2)
)

// GetLogger retrieves a logger instance
//
// Specify Conf.Modules during initialization to use the corresponding module
//
// Use the default module if the requested module has not been initialized
func GetLogger(module ...string) *zap.Logger {
	if len(module) < 1 {
		return gLoggerSvcHolders.load(gDefaultModeName).instance
	}

	return gLoggerSvcHolders.load(module[0]).instance
}

// Init initializes
//
// Shorthand for InitLoggerZap
func Init(cfg Conf, appName string, syncers ...zapcore.WriteSyncer) {
	InitLoggerZap(cfg, appName, syncers...)
}

// InitLoggerZap initializes the Zap logging service
//
// A default module namespace is added; when GetLogger is called without arguments,
// the default module namespace is used
//
// When the exporter is enabled, Exporter.Provider specifies the exporter implementation;
// an empty value disables the exporter.
// Provider currently supports:
//
// 'redis'、'redis-cluster'、'nats'、'kafka'
func InitLoggerZap(cfg Conf, appName string, syncers ...zapcore.WriteSyncer) {
	//Inject the default module namespace
	cfg.Modules = append(cfg.Modules, gDefaultModeName)
	sc := &svcHolders{}

	//Define the global logging component configuration
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(cfg.Level) {
	default:
		atomicLevel.SetLevel(zapcore.DebugLevel)
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
		if cfg.MaxSize == 0 {
			cfg.MaxSize = 100
		}
		if cfg.MaxBackups == 0 {
			cfg.MaxBackups = 10
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

		//Output to a file
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(fileWriter)), atomicLevel,
		)

		cores = append(cores, fileCore)

		//Output to the terminal (if enabled in configuration)
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

		//Define a single Logstash subscription key, independent of the module
		writer := exporterProvider(cfg)
		if writer != nil {
			coreWithWriter := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(writer)), atomicLevel,
			)
			cores = append(cores, coreWithWriter)
		}

		//Set a custom exporter
		if len(syncers) > 0 {
			gWriterSyncers = append(gWriterSyncers, syncers...)
		}

		//Read externally configured syncers and add them to the cores
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

// Set the exporter
func exporterProvider(cfg Conf) zapcore.WriteSyncer {
	var writer zapcore.WriteSyncer

	switch strings.ToLower(cfg.Exporter.Provider) {
	case "redis":
		redisInstance, err := redis.New(cfg.Exporter.Redis.ConnConf)
		if err != nil {
			log.Println("[logger] using (redis) exporter, but initialize connection error: ", err)
			return writer
		}
		redisWriter := &redisWriterStd{
			cli:     redisInstance,
			listKey: cfg.Exporter.Redis.ListKey,
		}

		writer = redisWriter
		log.Println("[logger] using (redis) exporter")
		return writer
	case "redis-cluster":
		redisInstance, err := redis.NewCluster(cfg.Exporter.Redis.ClusterConnConf)
		if err != nil {
			log.Println("[logger] using (redis-cluster) exporter, but initialize connection error: ", err)
			return writer
		}
		redisWriter := &redisClusterWriterStd{
			cli:     redisInstance,
			listKey: cfg.Exporter.Redis.ListKey,
		}

		writer = redisWriter
		log.Println("[logger] using (redis-cluster) exporter")
		return writer
	case "nats":
		natsInstance, err := nats.New(cfg.Exporter.Nats.ConnConf)
		if err != nil {
			log.Println("[logger] using (nats) exporter, but initialize connection error: ", err)
			return writer
		}
		natsWriter := &natsWriterStd{
			cli:        natsInstance,
			subjectKey: cfg.Exporter.Nats.Subject,
		}

		writer = natsWriter
		log.Println("[logger] using (nats) exporter")
		return writer
	case "kafka":
		kafkaInstance, err := kafka.NewWriter(cfg.Exporter.Kafka.ConnConf, cfg.Exporter.Kafka.Topic)
		if err != nil {
			log.Println("[logger] using (kafka) exporter, but initialize writer error: ", err)
			return writer
		}
		kafkaWriter := &kafkaWriterStd{
			writer: kafkaInstance,
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
