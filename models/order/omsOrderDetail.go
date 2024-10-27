package order

// OmsOrderDetail 包含商品信息的订单详情
type OmsOrderDetail struct {
	//订单信息（这个类继承自OmsOrder类）
	Order OmsOrder

	//订单商品列表
	OrderItemList []OmsOrderItem
}
