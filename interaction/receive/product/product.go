package product

type DetailReqStruct struct {
	Id int64 `json:"id"`
}

type SearchReqStruct struct {
	BrandId           int64  `json:"brandId"`
	Keyword           string `json:"keyword"`
	PageNum           int    `json:"pageNum"`
	PageSize          int    `json:"pageSize"`
	ProductCategoryId int64  `json:"productCategoryId"`
	Sort              int    `json:"sort"`
}
