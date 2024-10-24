package users

import (
	"gomall/global"
	"gomall/models/common"
	"time"
)

type User struct {
	common.PublicModel
	MemberLevelId         int64     `json:"member_level_id" gorm:"member_level_id"`
	Username              string    `json:"username" gorm:"username"`
	Password              string    `json:"password" gorm:"password"`
	Nickname              string    `json:"nickname" gorm:"nickname"`
	Phone                 string    `json:"phone" gorm:"phone"`
	Status                int       `json:"status" gorm:"status"`
	Icon                  string    `json:"icon" gorm:"icon"`
	Gender                int       `json:"gender" gorm:"gender"`
	Birthday              time.Time `json:"birthday" gorm:"birthday"`
	City                  string    `json:"city" gorm:"city"`
	Job                   string    `json:"job" gorm:"job"`
	PersonalizedSignature string    `json:"personalized_signature" gorm:"personalized_signature"`
	SourceType            int       `json:"source_type" gorm:"source_type"`
	Integration           int       `json:"integration" gorm:"integration"`
	Growth                int       `json:"growth" gorm:"growth"`
	//LuckyCount            int       `json:"lucky_count" gorm:"lucky_count"`
	HistoryIntegration int `json:"history_integration" gorm:"history_integration"`
}

func (User) TableName() string {
	return "ums_member"
}

func (us *User) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&us).Error
	if err != nil {
		return false
	}
	if us.ID <= 0 {
		return false
	}
	return true
}
