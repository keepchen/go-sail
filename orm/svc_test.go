package orm

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"gorm.io/gorm"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/utils"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
)

type User struct {
	BaseModel
	UserID   int64  `gorm:"column:user_id;type:bigint;not null;index:,unique;comment:用户ID"`
	Nickname string `gorm:"column:nickname;type:varchar(30);comment:用户昵称"`
	Status   int    `gorm:"column:status;type:tinyint;default:0;comment:用户状态"`
}

func (*User) TableName() string {
	return "user"
}

type Wallet struct {
	BaseModel
	UserID int64           `gorm:"column:user_id;type:bigint;not null;index:,unique;comment:用户ID"`
	Amount decimal.Decimal `gorm:"column:amount;type:decimal(10,2);default:0;comment:余额"`
}

func (*Wallet) TableName() string {
	return "wallet"
}

var (
	loggerConf = logger.Conf{
		Level:    "debug",
		Filename: "../examples/logs/testcase_db.log",
	}
	dbConf = db.Conf{
		Enable:      true,
		DriverName:  "mysql",
		AutoMigrate: true,
		Logger: db.Logger{
			Level: "debug",
		},
		ConnectionPool: db.ConnectionPoolConf{
			MaxOpenConnCount:       100,
			MaxIdleConnCount:       10,
			ConnMaxLifeTimeMinutes: 30,
			ConnMaxIdleTimeMinutes: 10,
		},
		Mysql: db.MysqlConf{
			Read: db.MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      3306,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
			Write: db.MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      3306,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
		},
	}
)

func TestSvcUsage(t *testing.T) {
	logger.Init(loggerConf, "go-sail")
	dbr, _, dbw, _ := db.New(dbConf)
	//logger.Init(loggerConf)
	if dbr == nil || dbw == nil {
		t.Log("database instance is nil, testing not emit.")
		return
	}
	_ = AutoMigrate(dbw, &User{}, &Wallet{})

	svc := New(dbr, dbw, logger.GetLogger())

	dbw.Exec(fmt.Sprintf("truncate table %s", (&User{}).TableName()))
	dbw.Exec(fmt.Sprintf("truncate table %s", (&Wallet{}).TableName()))

	SetHookTimeFunc(utils.NewTimeWithTimeZone(constants.TimeZoneUTCPlus7).Now)

	// ---- ignore gorm.ErrRecordNotFound
	var user0 User
	err := svc.R().FirstOrNil(&user0)
	t.Log("FirstOrNil:", err)
	assert.NoError(t, err)

	// ---- create record
	var user1 = User{
		UserID:   1000,
		Nickname: "go-sail",
		Status:   1,
	}
	err = svc.W().Create(&user1)
	assert.NoError(t, err)
	t.Log("Create user:", user1)

	var wallet1 = Wallet{
		UserID: user1.UserID,
		Amount: decimal.NewFromFloat(10.24),
	}

	err = svc.Create(&wallet1)
	assert.NoError(t, err)
	t.Log("Create wallet:", wallet1)

	// ---- read one record
	var user2 User
	err = svc.R().Where(&User{UserID: 1000}).First(&user2)
	assert.NoError(t, err)
	t.Log("First:", user2)

	// ---- force update all fields except some one
	var user3 = User{
		UserID:   1000,
		Nickname: "go-sail",
		Status:   2,
	}
	err = svc.W().Select("*").Omit("id", "created_at", "deleted_at").
		Where("user_id = ?", user3.UserID).Updates(&user3)
	assert.NoError(t, err)
	t.Log("Updates:", user3)

	var (
		queryUser   User
		queryWallet Wallet
	)

	err = svc.Model(&User{}).Where("user_id = ?", user3.UserID).First(&queryUser)
	assert.NoError(t, err)

	err = svc.Model(&Wallet{}).Where("amount = ?", 10.24).First(&queryWallet)
	assert.NoError(t, err)

	// ---- read several records
	var (
		users0 []User
	)
	err = svc.R().Where("id > ?", 0).Find(&users0)
	assert.NoError(t, err)
	t.Log("Find:", users0)

	// ---- ignore gorm.ErrRecordNotFound
	var users1 User
	err = svc.R().Where(&User{UserID: 99999}).FindOrNil(&users1)
	t.Log("FindOrNil:", err)
	assert.NoError(t, err)
	assert.Equal(t, users1.ID, NoneID)

	// ---- paginate
	var (
		users2   []User
		page     = 1
		pageSize = 50
	)
	total, err := svc.R().Paginate(&users2, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(users2)), total)
	t.Log("Paginate:", users2, total)
}

func TestOr(t *testing.T) {
	t.Run("Or", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Or("1=1")
	})
}

func TestWhere(t *testing.T) {
	t.Run("Where", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Where("1=1")
	})
}

func TestNot(t *testing.T) {
	t.Run("Not", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Not("1=1")
	})
}

func TestJoins(t *testing.T) {
	t.Run("Joins", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Joins("1=1")
	})
}

func TestOrder(t *testing.T) {
	t.Run("Order", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Order("1=1")
	})
}

func TestGroup(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Group("1=1")
	})
}

func TestHaving(t *testing.T) {
	t.Run("Having", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Having("1=1")
	})
}

func TestOffset(t *testing.T) {
	t.Run("Offset", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Offset(0)
	})
}

func TestLimit(t *testing.T) {
	t.Run("Limit", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Limit(100)
	})
}

func TestScopes(t *testing.T) {
	t.Run("Scopes", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Scopes(Paginate(1, 100))
	})
}

func TestSession(t *testing.T) {
	t.Run("Session", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Session(&gorm.Session{})
	})
}

func TestWithContext(t *testing.T) {
	t.Run("WithContext", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().WithContext(context.Background())
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("Unwrap", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		NewSvcImplSilent().Unwrap()
	})
}

func TestCount(t *testing.T) {
	t.Run("Count", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		var totalCount int64
		NewSvcImplSilent().Count(&totalCount)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		t.Log(NewSvcImplSilent().Create(&User{}))
	})
}

func TestSave(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		t.Log(NewSvcImplSilent().Save(&User{}))
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		t.Log(NewSvcImplSilent().Delete(&User{}))
	})
}

func TestTransaction(t *testing.T) {
	t.Run("Transaction()", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, _, dbw, _ := db.New(dbConf)
		//logger.Init(loggerConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		t.Log(NewSvcImplSilent().Transaction(func(tx *gorm.DB) error {
			tx.First(&User{})
			return nil
		}))
	})
}
