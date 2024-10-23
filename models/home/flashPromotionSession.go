package home

import "time"

type HomeFlashPromotionSession struct {
	Id         int64     `json:"id" gorm:"id"`
	Name       string    `json:"name" gorm:"name"`
	StartDate  string    `json:"start_time" gorm:"start_time"`
	EndDate    string    `json:"end_time" gorm:"end_time"`
	Status     int       `json:"status" gorm:"status"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
}

func (HomeFlashPromotionSession) TableName() string {
	return "sms_flash_promotion_session"
}
