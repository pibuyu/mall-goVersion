package home

// GetHomeContentRequestStruct 请求主页内容
type GetHomeContentRequestStruct struct{}

type PageHelper struct {
	PageNum  int `json:"pageNum" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type GetHotProductListRequestStruct struct {
	PageNum  int `json:"pageNum" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type GetNewProductListRequestStruct struct {
	PageNum  int `json:"pageNum" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type GetRecommendProductListRequestStruct struct {
	PageNum  int `form:"pageNum" binding:"required"`
	PageSize int `form:"pageSize" binding:"required"`
}

type GetSubjectListRequestStruct struct {
	CateId   int64 `json:"cateId"`
	PageNum  int   `json:"pageNum" binding:"required"`
	PageSize int   `json:"pageSize" binding:"required"`
}

type GetProductCateListRequestStruct struct {
	ParentId int64 `json:"parentId"`
}
