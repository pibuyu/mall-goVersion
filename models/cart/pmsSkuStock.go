package cart

import (
	"database/sql/driver"
	"encoding/json"
	"gomall/global"
)

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

type PmsSkuStockList []PmsSkuStock

func (list *PmsSkuStockList) GetByProductId(productId int64) (err error) {
	return global.Db.Model(&PmsSkuStock{}).Where("product_id = ?", productId).Find(list).Error
}
func (PmsSkuStock) TableName() string {
	return "pms_sku_stock"
}

// todo：这两个接口是豆包让实现的，不然会报错
// 为 PmsSkuStock 结构体实现 Valuer 接口，用于将结构体转换为数据库可存储的值
func (p PmsSkuStock) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// 为 PmsSkuStock 结构体实现 Scanner 接口，用于将从数据库读取的值转换为结构体
func (p *PmsSkuStock) Scanner(val interface{}) error {
	return json.Unmarshal(val.([]byte), &p)
}

func (stock *PmsSkuStock) GetSkuStockById(id int64) {
	global.Db.Model(&PmsSkuStock{}).Where("id = ?", id).Find(&stock)
}
