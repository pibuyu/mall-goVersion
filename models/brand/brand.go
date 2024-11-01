package brand

import "gomall/global"

type PmsBrand struct {
	ID                  int64  `json:"id" gorm:"column:id"`                                     // 主键ID
	Name                string `json:"name" gorm:"column:name"`                                 // 品牌名称
	FirstLetter         string `json:"firstLetter" gorm:"column:first_letter"`                  // 首字母
	Sort                int    `json:"sort" gorm:"column:sort"`                                 // 排序
	FactoryStatus       int    `json:"factoryStatus" gorm:"column:factory_status"`              // 是否为品牌制造商：0->不是；1->是
	ShowStatus          int    `json:"showStatus" gorm:"column:show_status"`                    // 显示状态
	ProductCount        int    `json:"productCount" gorm:"column:product_count"`                // 产品数量
	ProductCommentCount int    `json:"productCommentCount" gorm:"column:product_comment_count"` // 产品评论数量
	Logo                string `json:"logo" gorm:"column:logo"`                                 // 品牌logo
	BigPic              string `json:"bigPic" gorm:"column:big_pic"`                            // 专区大图
	BrandStory          string `json:"brandStory" gorm:"column:brand_story"`                    // 品牌故事
}

func (b *PmsBrand) GetById(brandId int64) error {
	return global.Db.Model(&PmsBrand{}).Where("id = ?", brandId).First(b).Error
}

func (PmsBrand) TableName() string {
	return "pms_brand"
}
