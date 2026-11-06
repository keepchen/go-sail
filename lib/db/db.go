package db

import (
	"time"

	"github.com/keepchen/go-sail/v3/lib/logger"

	"gorm.io/gorm"
)

type Instance struct {
	R *gorm.DB //Read instance
	W *gorm.DB //Write instance
}

// Database connection instance
var dbInstance *Instance

// InitDB initializes database connections
func InitDB(conf Conf) {
	dialectR, dialectW := conf.GenDialector()
	//read instance
	dbPtrR := mustInitDB(conf, dialectR)
	//write instance
	dbPtrW := mustInitDB(conf, dialectW)

	dbInstance = &Instance{
		R: dbPtrR,
		W: dbPtrW,
	}
}

// NewFreshDB instantiates fresh database connections
//
// rInstance is the read instance and wInstance is the write instance
func NewFreshDB(conf Conf) (rInstance *gorm.DB, rErr error, wInstance *gorm.DB, wErr error) {
	dialectR, dialectW := conf.GenDialector()
	rInstance, rErr = initDB(conf, dialectR)
	wInstance, wErr = initDB(conf, dialectW)

	return
}

func mustInitDB(conf Conf, dialect gorm.Dialector) *gorm.DB {
	loggerSvc := NewZapLoggerForGorm(logger.GetLogger(), conf)
	loggerSvc.SetAsDefault()

	gormConf := &gorm.Config{
		Logger:                                   loggerSvc,
		AllowGlobalUpdate:                        conf.AllowGlobalUpdate,
		SkipDefaultTransaction:                   conf.SkipDefaultTransaction,
		DisableNestedTransaction:                 conf.DisableNestedTransaction,
		DisableForeignKeyConstraintWhenMigrating: conf.DisableForeignKeyConstraintWhenMigrating,
	}

	if conf.NowFunc != nil {
		gormConf.NowFunc = conf.NowFunc
	}

	dbPtr, err := gorm.Open(dialect, gormConf)
	if err != nil {
		panic(err)
	}

	sqlDB, err := dbPtr.DB()
	if err != nil {
		panic(err)
	}

	if conf.ConnectionPool.MaxOpenConnCount == 0 {
		conf.ConnectionPool.MaxOpenConnCount = 100
	}
	if conf.ConnectionPool.MaxIdleConnCount == 0 {
		conf.ConnectionPool.MaxIdleConnCount = 10
	}
	if conf.ConnectionPool.ConnMaxLifeTimeMinutes == 0 {
		conf.ConnectionPool.ConnMaxLifeTimeMinutes = 30
	}
	if conf.ConnectionPool.ConnMaxIdleTimeMinutes == 0 {
		conf.ConnectionPool.ConnMaxIdleTimeMinutes = 10
	}

	sqlDB.SetMaxOpenConns(conf.ConnectionPool.MaxOpenConnCount)
	sqlDB.SetMaxIdleConns(conf.ConnectionPool.MaxIdleConnCount)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxLifeTimeMinutes))
	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxIdleTimeMinutes))

	return dbPtr
}

func initDB(conf Conf, dialect gorm.Dialector) (*gorm.DB, error) {
	loggerSvc := NewZapLoggerForGorm(logger.GetLogger(), conf)
	loggerSvc.SetAsDefault()
	dbPtr, err := gorm.Open(dialect, &gorm.Config{
		Logger:                                   loggerSvc,
		AllowGlobalUpdate:                        conf.AllowGlobalUpdate,
		SkipDefaultTransaction:                   conf.SkipDefaultTransaction,
		DisableNestedTransaction:                 conf.DisableNestedTransaction,
		DisableForeignKeyConstraintWhenMigrating: conf.DisableForeignKeyConstraintWhenMigrating,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := dbPtr.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(conf.ConnectionPool.MaxOpenConnCount)
	sqlDB.SetMaxIdleConns(conf.ConnectionPool.MaxIdleConnCount)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxLifeTimeMinutes))
	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(conf.ConnectionPool.ConnMaxIdleTimeMinutes))

	return dbPtr, nil
}

// GetInstance retrieves database instances
//
// Retrieve connections initialized by InitDB
func GetInstance() *Instance {
	return dbInstance
}

// New initializes fresh database connections
//
// rInstance is the read instance and wInstance is the write instance
func New(conf Conf) (rInstance *gorm.DB, rErr error, wInstance *gorm.DB, wErr error) {
	return NewFreshDB(conf)
}

// Init initializes database connections
//
// Shorthand for InitDB
func Init(conf Conf) {
	InitDB(conf)
}
