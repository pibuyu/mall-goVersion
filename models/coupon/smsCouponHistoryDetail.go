package coupon

// 优惠券领取历史详情（包括优惠券信息和关联关系）
type SmsCouponHistoryDetail struct {
	SmsCouponHistory
	//相关优惠券信息
	Coupon SmsCoupon `json:"coupon"`
	//优惠券关联的商品
	ProductRelationList []SmsCouponProductRelation `json:"productRelationList"`
	//优惠券关联的商品分类
	CategoryRelationList []SmsCouponProductCategoryRelation `json:"categoryRelationList"`
}
