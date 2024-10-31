package cart

import "gomall/global"

// 产品满减表(只针对同商品)
type PmsProductFullReduction struct {
	Id          int64   `json:"id" gorm:"id"`
	ProductId   int64   `json:"product_id" gorm:"product_id"`
	FullPrice   float32 `json:"full_price" gorm:"full_price"`
	ReducePrice float32 `json:"reduce_price" gorm:"reduce_price"`
}

type PmsProductFullReductionList []PmsProductFullReduction

func (list *PmsProductFullReductionList) GetByProductId(productId int64) (err error) {
	return global.Db.Model(&PmsProductFullReduction{}).Where("product_id = ?", productId).Find(list).Error
}

func (PmsProductFullReduction) TableName() string {
	return "pms_product_full_reduction"
}
