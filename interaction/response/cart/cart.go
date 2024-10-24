package cart

// GetCartProductResponseStruct获取购物车商品信息的返回体类型
type GetCartProductResponseStruct struct {
	//pms_product表
	Id   int64  `json:"id"`
	Name string `json:"name"`
	//todo:这个字段在返回的数据中为空值
	SubTitle                   string  `json:"sub_title" gorm:"column:sub_title"`
	Price                      float32 `json:"price"`
	Pic                        string  `json:"pic"`
	ProductAttributeCategoryId int64   `json:"product_attribute_category_id"`
	Stock                      int     `json:"stock"`
	//pa表
	AttrId   int64  `json:"attr_id"`
	AttrName string `json:"attr_name"`
	//ps表
	SkuId    int64   `json:"sku_id"`
	SkuCode  string  `json:"sku_code"`
	SkuPrice float32 `json:"sku_price"`
	SkuStock int     `json:"sku_stock"`
	SkuPic   string  `json:"sku_pic"`
}
