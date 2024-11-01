package readHistory

type CreateReadHistoryReqStruct struct {
	ProductId       int64   `json:"productId"`
	ProductName     string  `json:"productName"`
	ProductPic      string  `json:"productPic"`
	ProductPrice    float32 `json:"productPrice"`
	ProductSubTitle string  `json:"productSubTitle"`
}

type DeleteReadHistoryReqStruct struct {
	Ids []int64 `json:"ids"`
}
type ListReadHistoryReqStruct struct {
	PageNum  int64 `form:"pageNum"`
	PageSize int64 `form:"pageSize"`
}
