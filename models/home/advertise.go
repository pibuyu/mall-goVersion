package home

import "time"

type HomeAdvertise struct {
	Id         int64     `json:"id" gorm:"id"`
	Name       string    `json:"name" gorm:"name"`
	Type       int       `json:"type" gorm:"type"`
	Pic        string    `json:"pic" gorm:"pic"`
	StartTime  time.Time `json:"startTime" gorm:"start_time"`
	EndTime    time.Time `json:"endTime" gorm:"end_time"`
	Status     int       `json:"status" gorm:"status"`
	ClickCount int       `json:"clickCount" gorm:"click_count"`
	OrderCount int       `json:"orderCount" gorm:"order_count"`
	Url        string    `json:"url" gorm:"url"`
	Note       string    `json:"note" gorm:"note"`
	Sort       int       `json:"sort" gorm:"sort"`
}

type AdvertiseList []HomeAdvertise

func (HomeAdvertise) TableName() string {
	return "sms_home_advertise"
}
