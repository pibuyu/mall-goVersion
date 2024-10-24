package cart

// OmsCartItem 购物车表
type OmsCartItem struct {
	Id                int64   `json:"id" gorm:"id"`
	ProductId         int64   `json:"product_id" gorm:"product_id"`
	ProductSkuId      int64   `json:"product_sku_id" gorm:"product_sku_id"`
	MemberId          int64   `json:"member_id" gorm:"member_id"`
	Quantity          int     `json:"quantity" gorm:"quantity"`
	Price             float32 `json:"price" gorm:"price"`
	ProductPic        string  `json:"product_pic" gorm:"product_pic"`
	ProductName       string  `json:"product_name" gorm:"product_name"`
	ProductSubTitle   string  `json:"product_sub_title" gorm:"product_sub_title"`
	ProductSkuCode    string  `json:"product_sku_code" gorm:"product_sku_code"`
	MemberNickname    string  `json:"member_nickname" gorm:"member_nickname"`
	CreateDate        string  `json:"create_date" gorm:"create_date"`
	ModifyDate        string  `json:"modify_date" gorm:"modify_date"`
	DeleteStatus      int     `json:"delete_status" gorm:"delete_status"`
	ProductCategoryId int64   `json:"product_category_id" gorm:"product_category_id"`
	ProductBrand      string  `json:"product_brand" gorm:"product_brand"`
	ProductSn         string  `json:"product_sn" gorm:"product_sn"`
	ProductAttr       string  `json:"product_attr" gorm:"product_attr"`
}

type OmsCartItemList []OmsCartItem

func (OmsCartItem) TableName() string {
	return "oms_cart_item"
}
