package es

type DeleteReqStruct struct {
	Id int64 `form:"id"`
}

type CreateReqStruct struct {
	Id int64 `form:"id"`
}

type SimpleSearchReqStruct struct {
	Keyword  string `json:"keyword" binding:"required"`
	PageNum  int    `json:"pageNum" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
}
