package order

// OmsOrderDetail 包含商品信息的订单详情
type OmsOrderDetail struct {
	//订单信息（这个类继承自OmsOrder类）
	OmsOrder

	//订单商品列表
	//todo:This is a special comment.这里记得注释一下json格式，不注释的话返回值是大写的OrderItemList，前端读不到数据，排查了好久。。。
	OrderItemList []OmsOrderItem `json:"orderItemList"`
}
