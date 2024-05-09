//此ORM包的作用是提供了读写实例分离的调用机制，另外还包装一些常用的CRUD方法供开发者调用，
//注意，这个包的存在并而不是为了替代gorm。
//
//目前已经包装了常规的创建、查询、更新、删除和分页方法。并且接受传入外部logger，此间的操
//作方法日志会由传入的外部logger收集和输出。
//
//要指定读、写实例可以调用 R() 或者 W() 方法，请查阅orm_example.go文件。
//
//更高阶的方法调用，请使用gorm库提供的语法糖。

package service

import (
	"database/sql"
	"errors"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/orm/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ORMSvc interface {
	Model(value interface{}) ORMSvc
	Where(query interface{}, args ...interface{}) ORMSvc
	Or(query interface{}, args ...interface{}) ORMSvc
	Not(query interface{}, args ...interface{}) ORMSvc
	Joins(query string, args ...interface{}) ORMSvc
	Select(query interface{}, args ...interface{}) ORMSvc
	Omit(columns ...string) ORMSvc
	Order(value interface{}) ORMSvc
	Group(name string) ORMSvc
	Offset(offset int) ORMSvc
	Limit(limit int) ORMSvc
	Having(query interface{}, args ...interface{}) ORMSvc
	Scopes(fns ...func(*gorm.DB) *gorm.DB) ORMSvc

	Create(value interface{}) error
	Find(dest interface{}, conditions ...interface{}) error
	First(dest interface{}, conditions ...interface{}) error
	Updates(values interface{}) error
	Delete(value interface{}, conditions ...interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)

	//R 使用读实例
	R() ORMSvc
	//W 使用写实例
	W() ORMSvc
	//Paginate 分页查询（多行）
	//
	//参数:
	//
	//page 页码，从1开始
	//
	//pageSize 每页条数，默认为10
	//
	//返回:
	//
	//总条数和错误
	Paginate(dest interface{}, page, pageSize int) (int64, error)
}

type ORMSvcImpl struct {
	dbr    *gorm.DB
	dbw    *gorm.DB
	tx     *gorm.DB
	logger *zap.Logger
}

var _ ORMSvc = (*ORMSvcImpl)(nil)

var NewORMSvcImpl = func(dbr *gorm.DB, dbw *gorm.DB, logger *zap.Logger) ORMSvc {
	return &ORMSvcImpl{
		dbr:    dbr,
		dbw:    dbw,
		logger: logger,
	}
}

func (a *ORMSvcImpl) Model(value interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Model(value)

	return a
}

func (a *ORMSvcImpl) Where(query interface{}, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Where(query, args...)

	return a
}

func (a *ORMSvcImpl) Or(query interface{}, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Or(query, args...)

	return a
}

func (a *ORMSvcImpl) Not(query interface{}, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Not(query, args...)

	return a
}

func (a *ORMSvcImpl) Joins(query string, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Joins(query, args...)

	return a
}

func (a *ORMSvcImpl) Select(query interface{}, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Select(query, args...)

	return a
}

func (a *ORMSvcImpl) Omit(columns ...string) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Omit(columns...)

	return a
}

func (a *ORMSvcImpl) Order(value interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Order(value)

	return a
}

func (a *ORMSvcImpl) Group(name string) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Group(name)

	return a
}

func (a *ORMSvcImpl) Having(query interface{}, args ...interface{}) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Having(query, args...)

	return a
}

func (a *ORMSvcImpl) Offset(offset int) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Offset(offset)

	return a
}

func (a *ORMSvcImpl) Limit(limit int) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Limit(limit)

	return a
}

func (a *ORMSvcImpl) Scopes(fns ...func(*gorm.DB) *gorm.DB) ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Scopes(fns...)

	return a
}

func (a *ORMSvcImpl) Create(value interface{}) error {
	err := a.dbw.Create(value).Error

	if err != nil {
		a.logger.Error("[Database service]:Create:Error",
			zap.Any("value", logger.MarshalInterfaceValue(value)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *ORMSvcImpl) Find(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Find(&dest, conditions...).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error("[Database service]:Find:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *ORMSvcImpl) First(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.First(dest, conditions...).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error("[Database service]:First:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *ORMSvcImpl) Updates(values interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Updates(values).Error

	if err != nil {
		a.logger.Error("[Database service]:Updates:Error",
			zap.String("value", logger.MarshalInterfaceValue(values)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *ORMSvcImpl) Delete(value interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}
	err := a.tx.Delete(value, conditions...).Error

	if err != nil {
		a.logger.Error("[Database service]:Delete:Error",
			zap.String("value", logger.MarshalInterfaceValue(value)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *ORMSvcImpl) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	if a.tx == nil {
		a.tx = a.dbw
	}
	err = a.tx.Transaction(fc, opts...)

	if err != nil {
		a.logger.Error("[Database service]:Transaction:Error",
			zap.Errors("errors", []error{err}))
	}

	return
}

func (a *ORMSvcImpl) R() ORMSvc {
	if a.tx == nil {
		a.tx = a.dbr
	}

	return a
}

func (a *ORMSvcImpl) W() ORMSvc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	return a
}

func (a *ORMSvcImpl) Paginate(dest interface{}, page, pageSize int) (int64, error) {
	if a.tx == nil {
		a.tx = a.dbw
	}

	var total int64
	a.tx.Count(&total)

	err := a.tx.Scopes(model.Paginate(page, pageSize)).Find(&dest).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error("[Database service]:Paginate:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.Errors("errors", []error{err}))
	}

	return total, err
}
