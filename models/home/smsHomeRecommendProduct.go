package home

// SmsHomeRecommendProduct 热门商品推荐的表
type SmsHomeRecommendProduct struct {
	Id              int64  `json:"id" gorm:"id"`
	ProductId       int64  `json:"product_id" gorm:"product_id"`
	ProductName     string `json:"product_name" gorm:"product_name"`
	RecommendStatus int    `json:"recommend_status" gorm:"recommend_status"`
	Sort            int    `json:"sort" gorm:"sort"`
}

func (SmsHomeRecommendProduct) TableName() string {
	return "sms_home_recommend_product"
}
