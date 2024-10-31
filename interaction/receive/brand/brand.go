package brand

type DetailReqStruct struct {
	BrandId int64 `json:"brandId"`
}
type RecommendListReqStruct struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}
type ProductListReqStruct struct {
	BrandId  int64 `json:"brandId"`
	PageNum  int   `json:"pageNum"`
	PageSize int   `json:"pageSize"`
}
