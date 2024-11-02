package productCollection

import "time"

// AddReqStruct 收藏商品请求结构体。严格来说只需要MemberId和ProductId两个字段就行,其他字段都会在插入时赋值，这里尊重原项目，都写上。
type AddReqStruct struct {
	CreateTime      time.Time `json:"createTime"`
	Id              string    `json:"id"`
	MemberIcon      string    `json:"memberIcon"`
	MemberId        int64     `json:"memberId"`
	MemberNickname  string    `json:"memberNickname"`
	ProductId       int64     `json:"productId" binding:"required"`
	ProductName     string    `json:"productName"`
	ProductPic      string    `json:"productPic"`
	ProductPrice    float32   `json:"productPrice"`
	ProductSubTitle string    `json:"productSubTitle"`
}

type DeleteRespStruct struct {
	ProductId int64 `json:"productId" binding:"required"`
}

type DetailRespStruct struct {
	ProductId int64 `json:"productId"`
}

type ListReqStruct struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}
