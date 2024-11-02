package brand

type DetailReqStruct struct {
	BrandId int64 `json:"brandId"`
}
type RecommendListReqStruct struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}
type ProductListReqStruct struct {
	BrandId  int64 `form:"brandId"`
	PageNum  int   `form:"pageNum"`
	PageSize int   `form:"pageSize"`
}
