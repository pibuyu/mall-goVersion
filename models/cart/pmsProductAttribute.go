package cart

type PmsProductAttribute struct {
	Id                         int64  `json:"id" gorm:"id"`
	ProductAttributeCategoryId int64  `json:"product_attribute_category_id" gorm:"product_attribute_category_id"`
	Name                       string `json:"name" gorm:"name"`
	SelectType                 int    `json:"select_type" gorm:"select_type"`
	InputType                  int    `json:"input_type" gorm:"input_type"`
	InputList                  string `json:"input_list" gorm:"input_list"`
	Sort                       int    `json:"sort" gorm:"sort"`
	FilterType                 int    `json:"filter_type" gorm:"filter_type"`
	SearchType                 int    `json:"search_type" gorm:"search_type"`
	RelatedStatus              int    `json:"related_status" gorm:"related_status"`
	HandAddStatus              int    `json:"hand_add_status" gorm:"hand_add_status"`
	Type                       int    `json:"type" gorm:"type"`
}

func (PmsProductAttribute) TableName() string {
	return "pms_product_attribute"

}
