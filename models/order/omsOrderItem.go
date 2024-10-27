package order

// ID 订单明细 ID
type OmsOrderItem struct {
	ID int64 `json:"id" gorm:"column:id"`
	// OrderId 订单 id
	OrderId int64 `json:"orderId" gorm:"column:order_id"`
	// OrderSn 订单编号
	OrderSn string `json:"orderSn" gorm:"column:order_sn"`
	// ProductId 商品 id
	ProductId int64 `json:"productId" gorm:"column:product_id"`
	// ProductPic 商品图片
	ProductPic string `json:"productPic" gorm:"column:product_pic"`
	// ProductName 商品名称
	ProductName string `json:"productName" gorm:"column:product_name"`
	// ProductBrand 商品品牌
	ProductBrand string `json:"productBrand" gorm:"column:product_brand"`
	// ProductSn 商品编号
	ProductSn string `json:"productSn" gorm:"column:product_sn"`
	// ProductPrice 销售价格
	ProductPrice float32 `json:"productPrice" gorm:"column:product_price"`
	// ProductQuantity 购买数量
	ProductQuantity int `json:"productQuantity" gorm:"column:product_quantity"`
	// ProductSkuId 商品 sku 编号
	ProductSkuId int64 `json:"productSkuId" gorm:"column:product_sku_id"`
	// ProductSkuCode 商品 sku 条码
	ProductSkuCode string `json:"productSkuCode" gorm:"column:product_sku_code"`
	// ProductCategoryId 商品分类 id
	ProductCategoryId int64 `json:"productCategoryId" gorm:"column:product_category_id"`
	// PromotionName 商品促销名称
	PromotionName string `json:"promotionName" gorm:"column:promotion_name"`
	// PromotionAmount 商品促销分解金额
	PromotionAmount float32 `json:"promotionAmount" gorm:"column:promotion_amount"`
	// CouponAmount 优惠券优惠分解金额
	CouponAmount float32 `json:"couponAmount" gorm:"column:coupon_amount"`
	// IntegrationAmount 积分优惠分解金额
	IntegrationAmount float32 `json:"integrationAmount" gorm:"column:integration_amount"`
	// RealAmount 该商品经过优惠后的分解金额
	RealAmount float32 `json:"realAmount" gorm:"column:real_amount"`
	// GiftIntegration 商品积分
	GiftIntegration int `json:"giftIntegration" gorm:"column:gift_integration"`
	// GiftGrowth 商品成长值
	GiftGrowth int `json:"giftGrowth" gorm:"column:gift_growth"`
	// ProductAttr 商品销售属性
	ProductAttr string `json:"productAttr" gorm:"column:product_attr"`
}

func (OmsOrderItem) TableName() string {
	return "oms_order_item"

}
