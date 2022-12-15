package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/keepchen/go-sail/pkg/lib/redis"
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

func (sc svcHolders) load(modeName string) svcStd {
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
	defaultModeName  = "<defaultModeName>"
	loggerSvcHolders *svcHolders
)

//GetLogger 获取日志服务实例
func GetLogger(modeName ...string) *zap.Logger {
	if len(modeName) < 1 {
		return loggerSvcHolders.load(defaultModeName).instance
	}

	return loggerSvcHolders.load(modeName[0]).instance
}

//InitLoggerZap 初始化zap日志服务
//
//会加入默认的一个模块空间，当不传参调用GetLogger()时，
//就是使用默认的模块空间
//
//当启用elk时，logger使用redis队列作为媒介，需要在logstash侧配置对应的pipeline
//队列的key取决于日志文件名和appName的组合，如：
//日志文件名=logs/app.log，appName=app
//则，队列名称为=> app:logs/app.log
func InitLoggerZap(cfg Conf, appName string, modeName ...string) {
	//注入默认的空间模块
	modeName = append(modeName, defaultModeName)
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

	//logstash订阅的redis队列key只定义一个，与modeName无关
	var redisKeyForLogstash string
	if cfg.EnableELKWithRedisList {
		if len(cfg.RedisListKey) == 0 {
			redisKeyForLogstash = fmt.Sprintf("%s:%s", appName, cfg.Filename)
		} else {
			redisKeyForLogstash = cfg.RedisListKey
		}
		log.Println("[!] redis for elk setting is ENABLE, list key is:", redisKeyForLogstash)
	}

	for _, mn := range modeName {
		var (
			filename string
			cores    []zapcore.Core
		)
		if mn != defaultModeName {
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

		//启用基于redis list的elk日志写入
		if cfg.EnableELKWithRedisList {
			if redis.GetInstance() != nil {
				redisWriter := &redisWriterStd{
					listKey: redisKeyForLogstash,
					cli:     redis.GetInstance(),
				}
				redisCore := zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(redisWriter), atomicLevel,
				)

				cores = append(cores, redisCore)
			}
			if redis.GetClusterInstance() != nil {
				redisWriter := &redisClusterWriterStd{
					listKey: redisKeyForLogstash,
					cli:     redis.GetClusterInstance(),
				}
				redisCore := zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(redisWriter), atomicLevel,
				)

				cores = append(cores, redisCore)
			}
		}

		zapCore := zapcore.NewTee(cores...)

		loggerWithFields := zap.New(zapCore, zap.AddCaller()).With(zap.String("serviceName", fmt.Sprintf("%s:%s", appName, mn)))

		sc.store(mn, loggerWithFields)
	}

	loggerSvcHolders = sc

	defer func() {
		for _, mn := range modeName {
			_ = loggerSvcHolders.load(mn).instance.Sync()
		}
	}()
}
