package es

// 搜索商品的属性信息
type EsProductAttributeValue struct {
	ID int64 `json:"id"`

	//这两个字段是pms_product_attribute_value表的字段
	ProductAttributeId int64  `json:"productAttributeId"`
	Value              string `json:"value" elastic:"type:keyword"` // 属性值

	//后两个字段是pms_product_attribute表的字段
	Type int    `json:"type"`                        // 属性参数类型：0 -> 规格；1 -> 参数
	Name string `json:"name" elastic:"type:keyword"` // 属性名称
}
