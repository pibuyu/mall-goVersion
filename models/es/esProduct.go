package es

// 搜索商品的信息
type EsProduct struct {
	ID                  int64   `json:"id" bson:"_id,omitempty"` // @Id
	ProductSn           string  `json:"productSn" elastic:"type:keyword"`
	BrandId             int64   `json:"brandId"`
	BrandName           string  `json:"brandName" elastic:"type:keyword"`
	ProductCategoryId   int64   `json:"productCategoryId"`
	ProductCategoryName string  `json:"productCategoryName" elastic:"type:keyword"`
	Pic                 string  `json:"pic"`
	Name                string  `json:"name" elastic:"type:text,analyzer:ik_max_word"`
	SubTitle            string  `json:"subTitle" elastic:"type:text,analyzer:ik_max_word"`
	Keywords            string  `json:"keywords" elastic:"type:text,analyzer:ik_max_word"`
	Price               float32 `json:"price"`
	Sale                int     `json:"sale"`
	NewStatus           int     `json:"newStatus"`
	RecommandStatus     int     `json:"recommandStatus"`
	Stock               int     `json:"stock"`
	PromotionType       int     `json:"promotionType"`
	Sort                int     `json:"sort"`
	//上面的字段和pms_product的字段一一对应
	AttrValueList []EsProductAttributeValue `json:"attrValueList" elastic:"type:nested,fielddata:true"`
}
