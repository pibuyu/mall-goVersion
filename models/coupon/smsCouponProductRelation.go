package coupon

import "gomall/global"

// 优惠券关联的商品信息
type SmsCouponProductRelation struct {
	// ID 记录ID
	ID int64 `json:"id" gorm:"column:id"`
	// CouponId 优惠券ID
	CouponId int64 `json:"couponId" gorm:"column:coupon_id"`
	// ProductId 商品ID
	ProductId int64 `json:"productId" gorm:"column:product_id"`
	// ProductName 商品名称
	ProductName string `json:"productName" gorm:"column:product_name"`
	// ProductSn 商品编码
	ProductSn string `json:"productSn" gorm:"column:product_sn"`
}

type SmsCouponProductRelationList []SmsCouponProductRelation

func (list *SmsCouponProductRelationList) GetByProductId(productId int64) error {
	return global.Db.Model(&SmsCouponProductRelation{}).Where("product_id = ?", productId).Find(list).Error
}

func (SmsCouponProductRelation) TableName() string {
	return "sms_coupon_product_relation"
}
