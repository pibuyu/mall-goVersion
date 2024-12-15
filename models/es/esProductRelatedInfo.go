package es

// 搜索商品的关联信息，包括品牌名称，分类名称及属性
type EsProductRelatedInfo struct {
	BrandNames           []string      `json:"brandNames"`
	ProductCategoryNames []string      `json:"productCategoryNames"`
	ProductAttrs         []ProductAttr `json:"productAttrs"`
}

type ProductAttr struct {
	AttrId     int64    `json:"attrId"`
	AttrName   string   `json:"attrName"`
	AttrValues []string `json:"attrValues"`
}
