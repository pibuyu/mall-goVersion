package home

type Brand struct {
	Id                  int64  `json:"id" gorm:"id"`
	Name                string `json:"name" gorm:"name"`
	FirstLetter         string `json:"first_letter" gorm:"first_letter"`
	Sort                string `json:"sort" gorm:"sort"`
	FactoryStatus       int    `json:"factory_status" gorm:"factory_status"`
	ShowStatus          int    `json:"show_status" gorm:"show_status"`
	ProductCount        int    `json:"product_count" gorm:"product_count"`
	ProductCommentCount int    `json:"product_comment_count" gorm:"product_comment_count"`
	Logo                string `json:"logo" gorm:"logo"`
	BigPic              string `json:"big_pic" gorm:"big_pic"`
	BrandStory          string `json:"brand_story" gorm:"brand_story"`
}

type BrandList []Brand

func (Brand) TableName() string {
	return "pms_brand"
}
