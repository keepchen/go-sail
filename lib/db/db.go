package db

import (
	"time"

	"github.com/keepchen/go-sail/v3/lib/logger"

	"gorm.io/gorm"
)

type Instance struct {
	R *gorm.DB //读实例
	W *gorm.DB //写实例
}

// 数据库连接实例
var dbInstance *Instance

// InitDB 初始化数据库连接
func InitDB(conf Conf) {
	dialectR, dialectW := conf.GenDialector()
	//read instance
	dbPtrR := initDB(conf, dialectR)
	//write instance
	dbPtrW := initDB(conf, dialectW)

	dbInstance = &Instance{
		R: dbPtrR,
		W: dbPtrW,
	}
}

// NewFreshDB 实例化全新的数据库链接
//
// rInstance为读实例,wInstance为写实例
func NewFreshDB(conf Conf) (rInstance, wInstance *gorm.DB) {
	dialectR, dialectW := conf.GenDialector()
	rInstance, wInstance = initDB(conf, dialectR), initDB(conf, dialectW)

	return
}

func initDB(conf Conf, dialect gorm.Dialector) *gorm.DB {
	loggerSvc := NewZapLoggerForGorm(logger.GetLogger(), conf)
	loggerSvc.SetAsDefault()
	dbPtr, err := gorm.Open(dialect, &gorm.Config{
		Logger: loggerSvc,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := dbPtr.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(conf.ConnectionPool.MaxOpenConnCount)
	sqlDB.SetMaxIdleConns(conf.ConnectionPool.MaxIdleConnCount)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxLifeTimeMinutes))
	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxIdleTimeMinutes))

	return dbPtr
}

// GetInstance 获取数据库实例
//
// 获取由InitDB实例化后的连接
func GetInstance() *Instance {
	return dbInstance
}

// New 初始化化全新的数据库链接
//
// rInstance为读实例,wInstance为写实例
func New(conf Conf) (rInstance, wInstance *gorm.DB) {
	rInstance, wInstance = NewFreshDB(conf)

	return
}

// Init 初始化数据库连接
//
// InitDB 的语法糖
func Init(conf Conf) {
	InitDB(conf)
}
