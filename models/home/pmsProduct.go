package home

import (
	"database/sql/driver"
	"encoding/json"
	"gomall/global"
	"time"
)

type PmsProduct struct {
	Id                         int64     `json:"id" gorm:"id"`
	BrandId                    int64     `json:"brand_id" gorm:"brand_id"`
	ProductCategoryId          int64     `json:"product_category_id" gorm:"product_category_id"`
	FeightTemplateId           int64     `json:"feight_template_id" gorm:"feight_template_id"`
	ProductAttributeCategoryId int64     `json:"product_attribute_category_id" gorm:"product_attribute_category_id"`
	Name                       string    `json:"name" gorm:"name"`
	Pic                        string    `json:"pic" gorm:"pic"`
	ProductSn                  string    `json:"product_sn" gorm:"product_sn"`
	DeleteStatus               int       `json:"delete_status" gorm:"delete_status"`
	PublishStatus              int       `json:"publish_status" gorm:"publish_status"`
	NewStatus                  int       `json:"new_status" gorm:"new_status"`
	RecommandStatus            int       `json:"recommand_status" gorm:"recommand_status"`
	VerifyStatus               int       `json:"verify_status" gorm:"verify_status"`
	Sort                       int       `json:"sort" gorm:"sort"`
	Sale                       int       `json:"sale" gorm:"sale"`
	Price                      float32   `json:"price" gorm:"price"`
	PromotionPrice             float32   `json:"promotion_price" gorm:"promotion_price"`
	GiftGrowth                 int       `json:"gift_growth" gorm:"gift_growth"`
	GiftPoint                  int       `json:"gift_point" gorm:"gift_point"`
	UsePointLimit              int       `json:"use_point_limit" gorm:"use_point_limit"`
	SubTitle                   string    `json:"sub_title" gorm:"sub_title"`
	Description                string    `json:"description" gorm:"description"`
	OriginalPrice              float32   `json:"original_price" gorm:"original_price"`
	Stock                      int       `json:"stock" gorm:"stock"`
	LowStock                   int       `json:"low_stock" gorm:"low_stock"`
	Unit                       string    `json:"unit" gorm:"unit"`
	Weight                     float32   `json:"weight" gorm:"weight"`
	PreviewStatus              int       `json:"preview_status" gorm:"preview_status"`
	ServiceIds                 string    `json:"service_ids" gorm:"service_ids"`
	Keywords                   string    `json:"keywords" gorm:"keywords"`
	Note                       string    `json:"note" gorm:"note"`
	AlbumPics                  string    `json:"album_pics" gorm:"album_pics"`
	DetailTitle                string    `json:"detail_title" gorm:"detail_title"`
	DetailDesc                 string    `json:"detail_desc" gorm:"detail_desc"`
	DetailHtml                 string    `json:"detail_html" gorm:"detail_html"`
	DetailMobileHtml           string    `json:"detail_mobile_html" gorm:"detail_mobile_html"`
	PromotionStartTime         time.Time `json:"promotion_start_time" gorm:"promotion_start_time"`
	PromotionEndTime           time.Time `json:"promotion_end_time" gorm:"promotion_end_time"`
	PromotionPerLimit          int       `json:"promotion_per_limit" gorm:"promotion_per_limit"`
	PromotionType              int       `json:"promotion_type" gorm:"promotion_type"`
	BrandName                  string    `json:"brand_name" gorm:"brand_name"`
	ProductCategoryName        string    `json:"product_category_name" gorm:"product_category_name"`
}
type PmsProductList []PmsProduct

func (p *PmsProduct) GetById(id int64) (err error) {
	return global.Db.Model(&PmsProduct{}).Where("id", id).First(&p).Error
}

// 为PmsProduct结构体实现Valuer接口，用于将结构体转换为数据库可存储的值
func (p *PmsProduct) Value() (driver.Value, error) {
	// 将结构体转换为JSON字节数组
	return json.Marshal(p)
}

// 为PmsProduct结构体实现Scanner接口，用于将从数据库读取的值转换为结构体
func (p *PmsProduct) Scanner(val interface{}) error {
	// 从数据库读取的值通常是字节数组形式，将其转换为结构体
	return json.Unmarshal(val.([]byte), p)
}
func (PmsProduct) TableName() string {
	return "pms_product"
}
