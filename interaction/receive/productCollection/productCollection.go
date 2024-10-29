package productCollection

import "time"

type AddReqStruct struct {
	CreateTime      time.Time `json:"createTime"`
	Id              string    `json:"id"`
	MemberIcon      string    `json:"memberIcon"`
	MemberId        int64     `json:"memberId"`
	MemberNickname  string    `json:"memberNickname"`
	ProductId       int64     `json:"productId"`
	ProductName     string    `json:"productName"`
	ProductPic      string    `json:"productPic"`
	ProductPrice    float32   `json:"productPrice"`
	ProductSubTitle string    `json:"productSubTitle"`
}
