package users

import (
	"errors"
	"gomall/global"
	"time"
)

type User struct {
	Id                    int64     `json:"id" gorm:"id"`
	CreateTime            time.Time `json:"createTime" gorm:"create_time"`
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

func (u *User) GetByName(name string) error {
	return global.Db.Model(&User{}).Where("username = ?", name).Find(u).Error
}

func (user *User) GetMemberById(memberId int64) (err error) {
	if err := global.Db.Model(&User{}).Where("id = ?", memberId).First(&user).Error; err != nil {
		return errors.New("根据id查询用户出错:" + err.Error())
	}
	return
}

func (user *User) Update() (err error) {
	if err = global.Db.Model(&User{}).Where("id=?", user.Id).Updates(user).Error; err != nil {
		return errors.New("修改用户信息出错:" + err.Error())
	}
	return nil
}

func (us *User) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&us).Error
	if err != nil {
		return false
	}
	if us.Id <= 0 {
		return false
	}
	return true
}
