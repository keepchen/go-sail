package orders

import (
	"fmt"

	"github.com/keepchen/go-sail/v2/pkg/common/db/models"
	modelsEnum "github.com/keepchen/go-sail/v2/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v2/pkg/utils"
	"gorm.io/gorm"
)

type Order struct {
	models.BaseModel
	UserID        int64                      `gorm:"column:user_id;type:bigint;not null;index:user_id;comment:用户编号"`
	OrderNo       string                     `gorm:"column:order_no;type:varchar(40);not null;index:order_no,unique;comment:订单号"`
	ProductName   string                     `gorm:"column:product_name;type:varchar(255);not null;comment:商品名称"`
	ProductNum    int                        `gorm:"column:product_num;type:int;not null;comment:商品数量"`
	Status        modelsEnum.OrderStatusCode `gorm:"column:status;type:tinyint;not null;default:0;comment:订单状态"`
	Amount        float64                    `gorm:"column:amount;type:decimal(14,4);not null;comment:订单金额"`
	TransactionID string                     `gorm:"column:transaction_id;type:varchar(40);index:bet_transaction_id;comment:交易流水号"`
	Platform      int                        `gorm:"column:platform;type:tinyint;not null;default:0;comment:下单的客户端平台(0:web,1:app)"`
}

func (*Order) TableName() string {
	return "tb_orders"
}

func (*Order) Class() string {
	return "Order"
}

func (o *Order) BeforeCreate(_ *gorm.DB) (err error) {
	o.BaseModel.CreatedAt = utils.NewTimeWithTimeZone().Now()
	o.BaseModel.UpdatedAt = o.BaseModel.CreatedAt

	return nil
}

func (o *Order) BeforeUpdate(_ *gorm.DB) (err error) {
	o.BaseModel.UpdatedAt = utils.NewTimeWithTimeZone().Now()

	return nil
}

func (*Order) GenOrderNo() string {
	return fmt.Sprintf("%d%s", utils.NewTimeWithTimeZone().Now().UnixNano(), utils.RandomDigitalChars(5))
}
