package cart

import "gomall/global"

// 产品阶梯价格表(只针对同商品)
type PmsProductLadder struct {
	Id        int64   `json:"id" gorm:"id"`
	ProductId int64   `json:"product_id" gorm:"product_id"`
	Count     int     `json:"count" gorm:"count"`       //满足的商品数量
	Discount  float32 `json:"discount" gorm:"discount"` //折扣
	Price     float32 `json:"price" gorm:"price"`       //折后价格
}

type PmsProductLadderList []PmsProductLadder

func (list *PmsProductLadderList) GetByProductId(productId int64) (err error) {
	return global.Db.Model(&PmsProductLadder{}).Where("product_id = ?", productId).Find(list).Error
}
func (PmsProductLadder) TableName() string {
	return "pms_product_ladder"
}
