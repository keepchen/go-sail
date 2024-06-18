//此ORM包的作用是提供了读写实例分离的调用机制，另外还包装一些常用的CRUD方法供开发者调用，
//注意，这个包的存在并而不是为了替代gorm。
//
//目前已经包装了常规的创建、查询、更新、删除和分页方法。并且接受传入外部logger，此间的操
//作产生的日志会由传入的外部logger收集和输出。
//
//要指定读、写实例可以调用 R() 或者 W() 方法，请查阅svc_test.go文件。
//
//更高阶的方法调用，请使用gorm库提供的语法糖。

package orm

import (
	"context"
	"database/sql"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Svc interface {
	Model(value interface{}) Svc
	Where(query interface{}, args ...interface{}) Svc
	Or(query interface{}, args ...interface{}) Svc
	Not(query interface{}, args ...interface{}) Svc
	Joins(query string, args ...interface{}) Svc
	Select(query interface{}, args ...interface{}) Svc
	Omit(columns ...string) Svc
	Order(value interface{}) Svc
	Group(name string) Svc
	Offset(offset int) Svc
	Limit(limit int) Svc
	Having(query interface{}, args ...interface{}) Svc
	Scopes(fns ...func(*gorm.DB) *gorm.DB) Svc
	Session(session *gorm.Session) Svc
	WithContext(ctx context.Context) Svc

	Count(count *int64)
	Create(value interface{}) error
	Find(dest interface{}, conditions ...interface{}) error
	First(dest interface{}, conditions ...interface{}) error
	Updates(values interface{}) error
	Save(values interface{}) error
	Delete(value interface{}, conditions ...interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)

	// Unwrap 返回gorm原生实例
	Unwrap() *gorm.DB

	// R 使用读实例
	R() Svc
	// W 使用写实例
	W() Svc
	// Paginate 分页查询（多行）
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
	// FindOrNil 查询多条记录
	//
	//如果记录不存在忽略 gorm.ErrRecordNotFound 错误
	FindOrNil(dest interface{}, conditions ...interface{}) error
	// FirstOrNil 查询单条记录
	//
	//如果记录不存在忽略 gorm.ErrRecordNotFound 错误
	FirstOrNil(dest interface{}, conditions ...interface{}) error
}

type SvcImpl struct {
	dbr    *gorm.DB
	dbw    *gorm.DB
	tx     *gorm.DB
	logger *zap.Logger
}

var _ Svc = (*SvcImpl)(nil)

// New 初始化
//
// NewSvcImpl 的语法糖
var New = NewSvcImpl

// NewSvcImpl 初始化
var NewSvcImpl = func(dbr *gorm.DB, dbw *gorm.DB, logger *zap.Logger) Svc {
	return &SvcImpl{
		dbr:    dbr,
		dbw:    dbw,
		logger: logger,
	}
}

// NewSvcImplSilent 初始化
//
// 使用默认的读写实例和日志对象
var NewSvcImplSilent = func() Svc {
	return &SvcImpl{
		dbr:    db.GetInstance().R,
		dbw:    db.GetInstance().W,
		logger: logger.GetLogger(),
	}
}

func (a *SvcImpl) Model(value interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Model(value)

	return a
}

func (a *SvcImpl) Where(query interface{}, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Where(query, args...)

	return a
}

func (a *SvcImpl) Or(query interface{}, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Or(query, args...)

	return a
}

func (a *SvcImpl) Not(query interface{}, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Not(query, args...)

	return a
}

func (a *SvcImpl) Joins(query string, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Joins(query, args...)

	return a
}

func (a *SvcImpl) Select(query interface{}, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Select(query, args...)

	return a
}

func (a *SvcImpl) Omit(columns ...string) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Omit(columns...)

	return a
}

func (a *SvcImpl) Order(value interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Order(value)

	return a
}

func (a *SvcImpl) Group(name string) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Group(name)

	return a
}

func (a *SvcImpl) Having(query interface{}, args ...interface{}) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Having(query, args...)

	return a
}

func (a *SvcImpl) Offset(offset int) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Offset(offset)

	return a
}

func (a *SvcImpl) Limit(limit int) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Limit(limit)

	return a
}

func (a *SvcImpl) Scopes(fns ...func(*gorm.DB) *gorm.DB) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Scopes(fns...)

	return a
}

func (a *SvcImpl) Session(session *gorm.Session) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.Session(session)

	return a
}

func (a *SvcImpl) WithContext(ctx context.Context) Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx = a.tx.WithContext(ctx)

	return a
}

func (a *SvcImpl) Unwrap() *gorm.DB {
	if a.tx == nil {
		a.tx = a.dbw
	}

	return a.tx
}

func (a *SvcImpl) Count(count *int64) {
	if a.tx == nil {
		a.tx = a.dbw
	}

	a.tx.Count(count)
	a.clearTx()
}

func (a *SvcImpl) Create(value interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Create(value).Error
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Create:Error",
			zap.Any("value", logger.MarshalInterfaceValue(value)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) Find(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := IgnoreErrRecordNotFound(a.tx.Find(dest, conditions...))
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Find:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) FindOrNil(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := IgnoreErrRecordNotFound(a.tx.Find(dest, conditions...))
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:FindOrNil:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) First(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := IgnoreErrRecordNotFound(a.tx.First(dest, conditions...))
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:First:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) FirstOrNil(dest interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := IgnoreErrRecordNotFound(a.tx.First(dest, conditions...))
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:FirstOrNil:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) Updates(values interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Updates(values).Error
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Updates:Error",
			zap.String("value", logger.MarshalInterfaceValue(values)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) Save(values interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Save(values).Error
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Save:Error",
			zap.String("value", logger.MarshalInterfaceValue(values)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) Delete(value interface{}, conditions ...interface{}) error {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err := a.tx.Delete(value, conditions...).Error
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Delete:Error",
			zap.String("value", logger.MarshalInterfaceValue(value)),
			zap.String("conditions", logger.MarshalInterfaceValue(conditions)),
			zap.Errors("errors", []error{err}))
	}

	return err
}

func (a *SvcImpl) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	if a.tx == nil {
		a.tx = a.dbw
	}

	err = a.tx.Transaction(fc, opts...)
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Transaction:Error",
			zap.Errors("errors", []error{err}))
	}

	return
}

func (a *SvcImpl) R() Svc {
	if a.tx == nil {
		a.tx = a.dbr
	}

	return a
}

func (a *SvcImpl) W() Svc {
	if a.tx == nil {
		a.tx = a.dbw
	}

	return a
}

func (a *SvcImpl) Paginate(dest interface{}, page, pageSize int) (int64, error) {
	if a.tx == nil {
		a.tx = a.dbw
	}

	var count int64
	if a.tx.Statement.Model == nil {
		a.tx.Model(dest).Count(&count)
	} else {
		a.tx.Count(&count)
	}

	err := IgnoreErrRecordNotFound(a.tx.Scopes(Paginate(page, pageSize)).Find(dest))
	a.clearTx()

	if err != nil {
		a.logger.Error("[Database service]:Paginate:Error",
			zap.String("value", logger.MarshalInterfaceValue(dest)),
			zap.Errors("errors", []error{err}))
	}

	return count, err
}

func (a *SvcImpl) clearTx() {
	a.tx = nil
}
