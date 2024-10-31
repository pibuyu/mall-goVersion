package product

type PmsProductCategoryNode struct {
	PmsProductCategory
	Children []PmsProductCategoryNode `json:"children"` // 子分类集合
}
