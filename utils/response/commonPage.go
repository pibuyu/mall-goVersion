package response

import "math"

type CommonPage[T any] struct {
	PageNum   int   `json:"pageNum"`
	PageSize  int   `json:"pageSize"`
	TotalPage int   `json:"totalPage"`
	Total     int64 `json:"total"`
	List      []T   `json:"list"`
}

// 将返回值切片封装为分页结果
func ResetPage[T any](list []T, total int64, pageNum, pageSize int) CommonPage[T] {
	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))
	return CommonPage[T]{
		PageNum:   pageNum,
		PageSize:  pageSize,
		TotalPage: totalPage,
		Total:     total,
		List:      list,
	}
}
