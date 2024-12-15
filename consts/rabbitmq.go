package consts

type Queue struct {
	Exchange string
	Name     string
	RouteKey string
}

// 定义队列常量
var (
	//订单超时取消用的异步队列
	QUEUE_TTL_ORDER_CANCEL = Queue{"mall.order.direct.ttl", "mall.order.cancel.ttl", "mall.order.cancel.ttl"}

	//异步更新用户的积分和成长值用的队列
	UPDATE_USER_INFO = Queue{"mall.user.direct", "mall.user.update", "mall.user.update"}
)
