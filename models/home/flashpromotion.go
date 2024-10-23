package home

import "time"

type FlashPromotion struct {
	Id         int64     `json:"id" gorm:"id"`
	Title      string    `json:"title" gorm:"title"`
	StartDate  string    `json:"start_time" gorm:"start_time"`
	EndDate    string    `json:"end_time" gorm:"end_time"`
	Status     int       `json:"status" gorm:"status"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
}

func (FlashPromotion) TableName() string {
	return "sms_flash_promotion"
}
