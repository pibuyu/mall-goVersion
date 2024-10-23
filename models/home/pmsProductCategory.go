package home

type PmsProductCategory struct {
	Id           int64  `json:"id" gorm:"id"`
	ParentId     int64  `json:"parent_id" gorm:"parent_id"`
	Name         string `json:"name" gorm:"name"`
	Level        int    `json:"level" gorm:"level"`
	ProductCount int    `json:"product_count" gorm:"product_count"`
	ProductUnit  string `json:"product_unit" gorm:"product_unit"`
	NavStatus    int    `json:"nav_status" gorm:"nav_status"`
	ShowStatus   int    `json:"show_status" gorm:"show_status"`
	Sort         int    `json:"sort" gorm:"sort"`
	Icon         string `json:"icon" gorm:"icon"`
	Keywords     string `json:"keywords" gorm:"keywords"`
	Description  string `json:"description" gorm:"description"`
}
type ProductCategoryList []PmsProductCategory

func (PmsProductCategory) TableName() string {
	return "pms_product_category"
}
