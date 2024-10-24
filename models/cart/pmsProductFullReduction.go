package cart

// 产品满减表(只针对同商品)
type PmsProductFullReduction struct {
	Id          int64   `json:"id" gorm:"id"`
	ProductId   int64   `json:"product_id" gorm:"product_id"`
	FullPrice   float32 `json:"full_price" gorm:"full_price"`
	ReducePrice float32 `json:"reduce_price" gorm:"reduce_price"`
}

func (PmsProductFullReduction) TableName() string {
	return "pms_product_full_reduction"
}
