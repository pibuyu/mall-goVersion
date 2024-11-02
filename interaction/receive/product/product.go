package product

type DetailReqStruct struct {
	Id int64 `json:"id"`
}

type SearchReqStruct struct {
	BrandId           int64  `json:"brandId"`
	Keyword           string `json:"keyword"`
	PageNum           int    `form:"pageNum"`
	PageSize          int    `form:"pageSize"`
	ProductCategoryId int64  `form:"productCategoryId"`
	Sort              int    `form:"sort"`
}
