package home

type Brand struct {
	Id                  int64  `json:"id" gorm:"id"`
	Name                string `json:"name" gorm:"name"`
	FirstLetter         string `json:"firstLetter" gorm:"first_letter"`
	Sort                string `json:"sort" gorm:"sort"`
	FactoryStatus       int    `json:"factoryStatus" gorm:"factory_status"`
	ShowStatus          int    `json:"showStatus" gorm:"show_status"`
	ProductCount        int    `json:"productCount" gorm:"product_count"`
	ProductCommentCount int    `json:"productCommentCount" gorm:"product_comment_count"`
	Logo                string `json:"logo" gorm:"logo"`
	BigPic              string `json:"bigPic" gorm:"big_pic"`
	BrandStory          string `json:"brandStory" gorm:"brand_story"`
}

type BrandList []Brand

func (Brand) TableName() string {
	return "pms_brand"
}
