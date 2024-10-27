package order

import (
	"gomall/global"
	"time"
)

// OmsOrder 订单表
type OmsOrder struct {
	ID int64 `json:"id" gorm:"column:id"`
	// MemberID 会员 ID
	MemberID int64 `json:"memberId" gorm:"column:member_id"`
	// CouponID 优惠券 ID
	CouponID int64 `json:"couponId" gorm:"column:coupon_id"`
	// OrderSn 订单编号
	OrderSn string `json:"orderSn" gorm:"column:order_sn"`
	// CreateTime 提交时间
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	// MemberUsername 用户帐号
	MemberUsername string `json:"memberUsername" gorm:"column:member_username"`
	// TotalAmount 订单总金额
	TotalAmount float32 `json:"totalAmount" gorm:"column:total_amount"`
	// PayAmount 应付金额（实际支付金额）
	PayAmount float32 `json:"payAmount" gorm:"column:pay_amount"`
	// FreightAmount 运费金额
	FreightAmount float32 `json:"freightAmount" gorm:"column:freight_amount"`
	// PromotionAmount 促销优化金额（促销价、满减、阶梯价）
	PromotionAmount float32 `json:"promotionAmount" gorm:"column:promotion_amount"`
	// IntegrationAmount 积分抵扣金额
	IntegrationAmount float32 `json:"integrationAmount" gorm:"column:integration_amount"`
	// CouponAmount 优惠券抵扣金额
	CouponAmount float32 `json:"couponAmount" gorm:"column:coupon_amount"`
	// DiscountAmount 管理员后台调整订单使用的折扣金额
	DiscountAmount float32 `json:"discountAmount" gorm:"column:discount_amount"`
	// PayType 支付方式：0->未支付；1->支付宝；2->微信
	PayType int `json:"payType" gorm:"column:pay_type"`
	// SourceType 订单来源：0->PC订单；1->app订单
	SourceType int `json:"sourceType" gorm:"column:source_type"`
	// Status 订单状态：0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->无效订单
	Status int `json:"status" gorm:"column:status"`
	// OrderType 订单类型：0->正常订单；1->秒杀订单
	OrderType int `json:"orderType" gorm:"column:order_type"`
	// DeliveryCompany 物流公司(配送方式)
	DeliveryCompany string `json:"deliveryCompany" gorm:"column:delivery_company"`
	// DeliverySn 物流单号
	DeliverySn string `json:"deliverySn" gorm:"column:delivery_sn"`
	// AutoConfirmDay 自动确认时间（天）
	AutoConfirmDay int `json:"autoConfirmDay" gorm:"column:auto_confirm_day"`
	// Integration 可以获得的积分
	Integration int `json:"integration" gorm:"column:integration"`
	// Growth 可以活动的成长值
	Growth int `json:"growth" gorm:"column:growth"`
	// PromotionInfo 活动信息
	PromotionInfo string `json:"promotionInfo" gorm:"column:promotion_info"`
	// BillType 发票类型：0->不开发票；1->电子发票；2->纸质发票
	BillType int `json:"billType" gorm:"column:bill_type"`
	// BillHeader 发票抬头
	BillHeader string `json:"billHeader" gorm:"column:bill_header"`
	// BillContent 发票内容
	BillContent string `json:"billContent" gorm:"column:bill_content"`
	// BillReceiverPhone 收票人电话
	BillReceiverPhone string `json:"billReceiverPhone" gorm:"column:bill_receiver_phone"`
	// BillReceiverEmail 收票人邮箱
	BillReceiverEmail string `json:"billReceiverEmail" gorm:"column:bill_receiver_email"`
	// ReceiverName 收货人姓名
	ReceiverName string `json:"receiverName" gorm:"column:receiver_name"`
	// ReceiverPhone 收货人电话
	ReceiverPhone string `json:"receiverPhone" gorm:"column:receiver_phone"`
	// ReceiverPostCode 收货人邮编
	ReceiverPostCode string `json:"receiverPostCode" gorm:"column:receiver_post_code"`
	// ReceiverProvince 省份/直辖市
	ReceiverProvince string `json:"receiverProvince" gorm:"column:receiver_province"`
	// ReceiverCity 城市
	ReceiverCity string `json:"receiverCity" gorm:"column:receiver_city"`
	// ReceiverRegion 区
	ReceiverRegion string `json:"receiverRegion" gorm:"column:receiver_region"`
	// ReceiverDetailAddress 详细地址
	ReceiverDetailAddress string `json:"receiverDetailAddress" gorm:"column:receiver_detail_address"`
	// Note 订单备注
	Note string `json:"note" gorm:"column:note"`
	// ConfirmStatus 确认收货状态：0->未确认；1->已确认
	ConfirmStatus int `json:"confirmStatus" gorm:"column:confirm_status"`
	// DeleteStatus 删除状态：0->未删除；1->已删除
	DeleteStatus int `json:"deleteStatus" gorm:"column:delete_status"`
	// UseIntegration 下单时使用的积分
	UseIntegration int `json:"useIntegration" gorm:"column:use_integration"`
	// PaymentTime 支付时间
	PaymentTime *time.Time `json:"paymentTime" gorm:"column:payment_time"`
	// DeliveryTime 发货时间
	DeliveryTime *time.Time `json:"deliveryTime" gorm:"column:delivery_time"`
	// ReceiveTime 确认收货时间
	ReceiveTime *time.Time `json:"receiveTime" gorm:"column:receive_time"`
	// CommentTime 评价时间
	CommentTime *time.Time `json:"commentTime" gorm:"column:comment_time"`
	// ModifyTime 修改时间
	ModifyTime *time.Time `json:"modifyTime" gorm:"column:modify_time"`
}

func (order *OmsOrder) Insert() error {
	return global.Db.Create(order).Error
}

func (OmsOrder) TableName() string {
	return "oms_order"
}
