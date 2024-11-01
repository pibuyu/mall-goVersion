package cart

import "gomall/global"

// OmsCartItem 购物车表
type OmsCartItem struct {
	Id                int64   `json:"id" gorm:"id"`
	ProductId         int64   `json:"productId" gorm:"product_id"`
	ProductSkuId      int64   `json:"productSkuId" gorm:"product_sku_id"`
	MemberId          int64   `json:"memberId" gorm:"member_id"`
	Quantity          int     `json:"quantity" gorm:"quantity"`
	Price             float32 `json:"price" gorm:"price"`
	ProductPic        string  `json:"productPic" gorm:"product_pic"`
	ProductName       string  `json:"productName" gorm:"product_name"`
	ProductSubTitle   string  `json:"productSubTitle" gorm:"product_sub_title"`
	ProductSkuCode    string  `json:"productSkuCode" gorm:"product_sku_code"`
	MemberNickname    string  `json:"memberNickname" gorm:"member_nickname"`
	CreateDate        string  `json:"createDate" gorm:"create_date"`
	ModifyDate        string  `json:"modifyDate" gorm:"modify_date"`
	DeleteStatus      int     `json:"deleteStatus" gorm:"delete_status"`
	ProductCategoryId int64   `json:"productCategoryId" gorm:"product_category_id"`
	ProductBrand      string  `json:"productBrand" gorm:"product_brand"`
	ProductSn         string  `json:"productSn" gorm:"product_sn"`
	ProductAttr       string  `json:"productAttr" gorm:"product_attr"`
}

func (item *OmsCartItem) GetById(id int64) error {
	return global.Db.Model(&OmsCartItem{}).Where("id = ?", id).Find(item).Error

}

type OmsCartItemList []OmsCartItem

func (OmsCartItem) TableName() string {
	return "oms_cart_item"
}
