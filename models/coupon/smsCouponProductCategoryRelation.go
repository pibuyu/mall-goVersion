package coupon

// 优惠券关联的商品所属的分类
type SmsCouponProductCategoryRelation struct {
	// ID 记录ID
	ID int64 `json:"id" gorm:"column:id"`
	// CouponId 优惠券ID
	CouponId int64 `json:"couponId" gorm:"column:coupon_id"`
	// ProductCategoryId 产品分类ID
	ProductCategoryId int64 `json:"productCategoryId" gorm:"column:product_category_id"`
	// ProductCategoryName 产品分类名称
	ProductCategoryName string `json:"productCategoryName" gorm:"column:product_category_name"`
	// ParentCategoryName 父分类名称
	ParentCategoryName string `json:"parentCategoryName" gorm:"column:parent_category_name"`
}

func (SmsCouponProductCategoryRelation) TableName() string {
	return "sms_coupon_product_category_relation"

}
