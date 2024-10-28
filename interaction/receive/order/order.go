package order

type DetailReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
}

type CancelOrderReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
}

type ConfirmReceiveOrderReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
}

type DeleteOrderReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
}

type GenerateConfirmOrderReqStruct struct {
	CartIds []int64 `json:"cartIds" binding:"required"`
}

type GenerateOrderReqStruct struct {
	CartIds                []int64 `json:"cartIds"`
	CouponId               int64   `json:"couponId"`
	MemberReceiveAddressId int64   `json:"memberReceiveAddressId"`
	PayType                int     `json:"payType"`
	UseIntegration         int     `json:"useIntegration"`
}

type ListReqStruct struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	Status   int `json:"status" binding:"required"`
}

type PaySuccessReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
	PayType int   `json:"payType" binding:"required"`
}
