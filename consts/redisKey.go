package consts

const (

	//用户token后缀
	TokenString = "tokenString"

	//用户信息的前缀
	RedisDatabase = "MALL"
	// REDIS_KEY_AUTH_CODE 验证码前缀
	REDIS_KEY_AUTH_CODE = "AUTH_CODE"

	REDIS_KEY_ORDER_ID = "oms:orderId"

	//用户下单相关
	AVOID_REPEAT_ORDER_PREFIX   = "AVOID_REPEAT_ORDER_PREFIX_"   //避免重复生成订单
	AVOID_REPEAT_PAYMENT_PREFIX = "AVOID_REPEAT_PAYMENT_PREFIX_" //避免重复支付

	Order2MemberIdMap = "Order2MemberIdMap"

	//生成订单相关
	INVENTORY_PRODUCT_SEGMENT = "INVENTORY_PRODUCT_SEGMENT_"
)
