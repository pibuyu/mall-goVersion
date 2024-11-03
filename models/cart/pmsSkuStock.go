package cart

import (
	"database/sql/driver"
	"encoding/json"
	"gomall/global"
)

type PmsSkuStock struct {
	Id             int64   `gorm:"column:id" json:"id"`                          // SKU ID
	ProductId      int64   `gorm:"column:product_id" json:"productId"`           // 产品ID
	SkuCode        string  `gorm:"column:sku_code" json:"skuCode"`               // SKU编码
	Price          float32 `gorm:"column:price" json:"price"`                    // 价格
	Stock          int     `gorm:"column:stock" json:"stock"`                    // 库存
	LowStock       int     `gorm:"column:low_stock" json:"lowStock"`             // 低库存警告
	Pic            string  `gorm:"column:pic" json:"pic"`                        // 图片
	Sale           int     `gorm:"column:sale" json:"sale"`                      // 销售量
	PromotionPrice float32 `gorm:"column:promotion_price" json:"promotionPrice"` // 促销价格
	LockStock      int     `gorm:"column:lock_stock" json:"lockStock"`           // 锁定库存
	SpData         string  `gorm:"column:sp_data" json:"spData"`                 // SKU数据
}

type PmsSkuStockList []PmsSkuStock

func (list *PmsSkuStockList) GetByProductId(productId int64) (err error) {
	return global.Db.Model(&PmsSkuStock{}).Where("product_id = ?", productId).Find(list).Error
}
func (PmsSkuStock) TableName() string {
	return "pms_sku_stock"
}

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
