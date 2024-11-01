package home

import (
	"database/sql/driver"
	"encoding/json"
	"gomall/global"
	"time"
)

type PmsProduct struct {
	Id                         int64     `json:"id" gorm:"id"`
	BrandId                    int64     `json:"brandId" gorm:"brand_id"`
	ProductCategoryId          int64     `json:"productCategoryId" gorm:"product_category_id"`
	FeightTemplateId           int64     `json:"feightTemplateId" gorm:"feight_template_id"`
	ProductAttributeCategoryId int64     `json:"productAttributeCategoryId" gorm:"product_attribute_category_id"`
	Name                       string    `json:"name" gorm:"name"`
	Pic                        string    `json:"pic" gorm:"pic"`
	ProductSn                  string    `json:"productSn" gorm:"product_sn"`
	DeleteStatus               int       `json:"deleteStatus" gorm:"delete_status"`
	PublishStatus              int       `json:"publishStatus" gorm:"publish_status"`
	NewStatus                  int       `json:"newStatus" gorm:"new_status"`
	RecommandStatus            int       `json:"recommandStatus" gorm:"recommand_status"`
	VerifyStatus               int       `json:"verifyStatus" gorm:"verify_status"`
	Sort                       int       `json:"sort" gorm:"sort"`
	Sale                       int       `json:"sale" gorm:"sale"`
	Price                      float32   `json:"price" gorm:"price"`
	PromotionPrice             float32   `json:"promotionPrice" gorm:"promotion_price"`
	GiftGrowth                 int       `json:"giftGrowth" gorm:"gift_growth"`
	GiftPoint                  int       `json:"giftPoint" gorm:"gift_point"`
	UsePointLimit              int       `json:"usePointLimit" gorm:"use_point_limit"`
	SubTitle                   string    `json:"subTitle" gorm:"sub_title"`
	Description                string    `json:"description" gorm:"description"`
	OriginalPrice              float32   `json:"originalPrice" gorm:"original_price"`
	Stock                      int       `json:"stock" gorm:"stock"`
	LowStock                   int       `json:"lowStock" gorm:"low_stock"`
	Unit                       string    `json:"unit" gorm:"unit"`
	Weight                     float32   `json:"weight" gorm:"weight"`
	PreviewStatus              int       `json:"previewStatus" gorm:"preview_status"`
	ServiceIds                 string    `json:"serviceIds" gorm:"service_ids"`
	Keywords                   string    `json:"keywords" gorm:"keywords"`
	Note                       string    `json:"note" gorm:"note"`
	AlbumPics                  string    `json:"albumPics" gorm:"album_pics"`
	DetailTitle                string    `json:"detailTitle" gorm:"detail_title"`
	DetailDesc                 string    `json:"detailDesc" gorm:"detail_desc"`
	DetailHtml                 string    `json:"detailHtml" gorm:"detail_html"`
	DetailMobileHtml           string    `json:"detailMobileHtml" gorm:"detail_mobile_html"`
	PromotionStartTime         time.Time `json:"promotionStartTime" gorm:"promotion_start_time"`
	PromotionEndTime           time.Time `json:"promotionEndTime" gorm:"promotion_end_time"`
	PromotionPerLimit          int       `json:"promotionPerLimit" gorm:"promotion_per_limit"`
	PromotionType              int       `json:"promotionType" gorm:"promotion_type"`
	BrandName                  string    `json:"brandName" gorm:"brand_name"`
	ProductCategoryName        string    `json:"productCategoryName" gorm:"product_category_name"`
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
