package brand

import "gomall/global"

type PmsBrand struct {
	ID                  int64  `gorm:"column:id" json:"id"`                                     // 主键ID
	Name                string `gorm:"column:name" json:"name"`                                 // 品牌名称
	FirstLetter         string `gorm:"column:first_letter" json:"firstLetter"`                  // 首字母
	Sort                int    `gorm:"column:sort" json:"sort"`                                 // 排序
	FactoryStatus       int    `gorm:"column:factory_status" json:"factoryStatus"`              // 是否为品牌制造商：0->不是；1->是
	ShowStatus          int    `gorm:"column:show_status" json:"showStatus"`                    // 显示状态
	ProductCount        int    `gorm:"column:product_count" json:"productCount"`                // 产品数量
	ProductCommentCount int    `gorm:"column:product_comment_count" json:"productCommentCount"` // 产品评论数量
	Logo                string `gorm:"column:logo" json:"logo"`                                 // 品牌logo
	BigPic              string `gorm:"column:big_pic" json:"bigPic"`                            // 专区大图
	BrandStory          string `gorm:"column:brand_story" json:"brandStory"`                    // 品牌故事
}

func (b *PmsBrand) GetById(brandId int64) error {
	return global.Db.Model(&PmsBrand{}).Where("id = ?", brandId).First(b).Error
}

func (PmsBrand) TableName() string {
	return "pms_brand"
}
