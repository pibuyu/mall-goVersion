package consts

type Queue struct {
	Exchange string
	Name     string
	RouteKey string
}

// 定义队列常量
var (
	QUEUE_ORDER_CANCEL       = Queue{"mall.order.direct", "mall.order.cancel", "mall.order.cancel"}
	QUEUE_TTL_ORDER_CANCEL   = Queue{"mall.order.direct.ttl", "mall.order.cancel.ttl", "mall.order.cancel.ttl"}
	DELAY_QUEUE_ORDER_CANCEL = Queue{"mall.order.delay.direct", "mall.order.cancel.delay", "mall.order.cancel.delay"}
)
