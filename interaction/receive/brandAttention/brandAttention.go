package brandAttention

import "time"

type AddBrandAttentionReqStruct struct {
	BrandCity      string    `json:"brandCity"`
	BrandId        int64     `json:"brandId"`
	BrandLogo      string    `json:"brandLogo"`
	BrandName      string    `json:"brandName"`
	CreateTime     time.Time `json:"createTime"`
	Id             string    `json:"id"`
	MemberIcon     string    `json:"memberIcon"`
	MemberId       int64     `json:"memberId"`
	MemberNickname string    `json:"memberNickname"`
}

type DeleteBrandAttentionReqStruct struct {
	BrandId int64 `json:"brandId"`
}

type DetailReqStruct struct {
	BrandId int64 `json:"brandId"`
}

type ListReqStruct struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}
