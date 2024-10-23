package home

// FlashPromotionProductRelation 商品与秒杀活动的关系表
type FlashPromotionProductRelation struct {
	Id                      int64   `json:"id" gorm:"id"`
	FlashPromotionId        int64   `json:"flash_promotion_id" gorm:"flash_promotion_id"`
	FlashPromotionSessionId int64   `json:"flash_promotion_session_id" gorm:"flash_promotion_session_id"`
	ProductId               int64   `json:"product_id" gorm:"product_id"`
	FlashPromotionPrice     float32 `json:"flash_promotion_price" gorm:"flash_promotion_price"`
	FlashPromotionCount     int     `json:"flash_promotion_count" gorm:"flash_promotion_count"`
	FlashPromotionLimit     int     `json:"flash_promotion_limit" gorm:"flash_promotion_limit"`
	Sort                    int     `json:"sort" gorm:"sort"`
}

func (FlashPromotionProductRelation) TableName() string {
	return "sms_flash_promotion_product_relation"
}
