package models

type OrderStatusCode int

const (
	OrderStatusCodeCreated   = OrderStatusCode(0) //Order created
	OrderStatusCodePaid      = OrderStatusCode(1) //Order paid
	OrderStatusCodeClosed    = OrderStatusCode(2) //Order closed
	OrderStatusCodeCancelled = OrderStatusCode(3) //Order canceled
	OrderStatusCodeRefunded  = OrderStatusCode(4) //Order refunded
)

var orderStatusCodeMap = map[OrderStatusCode]string{
	OrderStatusCodeCreated:   "订单已创建",
	OrderStatusCodePaid:      "订单已支付",
	OrderStatusCodeClosed:    "订单已关闭",
	OrderStatusCodeCancelled: "订单已取消",
	OrderStatusCodeRefunded:  "订单已退款",
}

func (osc OrderStatusCode) Int() int {
	return int(osc)
}

func (osc OrderStatusCode) String() string {
	return orderStatusCodeMap[osc]
}
