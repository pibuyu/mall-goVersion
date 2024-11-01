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

// 该请求形如：http://localhost:9090/order/list?status=-1&pageNum=1&pageSize=5，需要在请求的结构体中用form注释
type ListReqStruct struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
	//传递过来的status=-1，不能用binding=required注释
	Status int `form:"status"`
}

type PaySuccessReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
	PayType int   `json:"payType" binding:"required"`
}

type CancelUserOrderReqStruct struct {
	OrderId int64 `json:"orderId" binding:"required"`
}

type CreateReturnApplyReqStruct struct {
	OrderId          int64   `json:"orderId"`
	ProductId        int64   `json:"productId"`
	OrderSn          string  `json:"orderSn"`
	MemberUsername   string  `json:"memberUsername"`
	ReturnName       string  `json:"returnName"`
	ReturnPhone      string  `json:"returnPhone"`
	ProductPic       string  `json:"productPic"`
	ProductName      string  `json:"productName"`
	ProductBrand     string  `json:"productBrand"`
	ProductAttr      string  `json:"productAttr"`
	ProductCount     int     `json:"productCount"`
	ProductPrice     float64 `json:"productPrice"`
	ProductRealPrice float64 `json:"productRealPrice"`
	Reason           string  `json:"reason"`
	Description      string  `json:"description"`
	ProofPics        string  `json:"proofPics"`
}
