package coupon

type AddCouponRequestStruct struct {
	CouponId int64 `json:"couponId" binding:"required"`
}

// todo:This is a special comment. 传递过来的useStatus=0也是有效值，当指定binding=required时会将零值视为没有传参，引发错误。因此此处去掉binding=required。
type ListCouponRequestStruct struct {
	UseStatus int `form:"useStatus"`
}

type ListCartRequestStruct struct {
	Type int `json:"type"`
}

type ListByProductRequestStruct struct {
	ProductId int64 `json:"productId"`
}

type ListHistoryRequestStruct struct {
	UseStatus int `json:"useStatus"`
}
