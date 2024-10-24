package cart

type PmsSkuStock struct {
	Id             int64   `json:"id" gorm:"id"`
	ProductId      int64   `json:"product_id" gorm:"product_id"`
	SkuCode        string  `json:"sku_code" gorm:"sku_code"`
	Price          float32 `json:"price" gorm:"price"`
	Stock          int     `json:"stock" gorm:"stock"`
	LowStock       int     `json:"low_stock" gorm:"low_stock"`
	Pic            string  `json:"pic" gorm:"pic"`
	Sale           int     `json:"sale" gorm:"sale"`
	PromotionPrice float32 `json:"promotion_price" gorm:"promotion_price"`
	LockStock      int     `json:"lock_stock" gorm:"lock_stock"`
	SpData         string  `json:"sp_data" gorm:"sp_data"`
}

func (PmsSkuStock) TableName() string {
	return "pms_sku_stock"
}
