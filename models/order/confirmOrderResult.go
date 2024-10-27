package order

import (
	"gomall/models/cart"
	"gomall/models/coupon"
	"gomall/models/integration"
)

type ConfirmOrderResult struct {
	// CartPromotionItemList 包含优惠信息的购物车信息列表
	CartPromotionItemList []cart.CartPromotionItem `json:"cartPromotionItemList"`
	// MemberReceiveAddressList 用户收货地址列表
	MemberReceiveAddressList []UmsMemberReceiveAddress `json:"memberReceiveAddressList"`
	// CouponHistoryDetailList 用户可用优惠券列表
	CouponHistoryDetailList []coupon.SmsCouponHistoryDetail `json:"couponHistoryDetailList"`
	// IntegrationConsumeSetting 积分使用规则
	IntegrationConsumeSetting integration.UmsIntegrationConsumeSetting `json:"integrationConsumeSetting"`
	// MemberIntegration 会员持有的积分
	MemberIntegration int `json:"memberIntegration"`
	// CalcAmount 计算的金额
	CalcAmount CalcAmount `json:"calcAmount"`
}

type CalcAmount struct {
	// TotalAmount 订单商品总金额
	TotalAmount float32 `json:"totalAmount" gorm:"column:total_amount"`
	// FreightAmount 运费
	FreightAmount float32 `json:"freightAmount" gorm:"column:freight_amount"`
	// PromotionAmount 活动优惠金额
	PromotionAmount float32 `json:"promotionAmount" gorm:"column:promotion_amount"`
	// PayAmount 应付金额
	PayAmount float32 `json:"payAmount" gorm:"column:pay_amount"`
}
