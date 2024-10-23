package home

type HomeFlashPromotionProduct struct {
	Product             PmsProduct
	FlashPromotionPrice float32 `json:"flash_promotion_price"`
	FlashPromotionCount int     `json:"flash_promotion_count"`
	FlashPromotionLimit int     `json:"flash_promotion_limit"`
}
type FlashPromotionProductList []HomeFlashPromotionProduct
