package order

import "gomall/global"

type UmsMemberReceiveAddress struct {
	// ID 地址ID
	ID int64 `json:"id" gorm:"column:id"`
	// MemberId 会员ID
	MemberId int64 `json:"memberId" gorm:"column:member_id"`
	// Name 收货人名称
	Name string `json:"name" gorm:"column:name"`
	// PhoneNumber 电话号码
	PhoneNumber string `json:"phoneNumber" gorm:"column:phone_number"`
	// DefaultStatus 是否为默认地址状态，0否，1是
	DefaultStatus int `json:"defaultStatus" gorm:"column:default_status"`
	// PostCode 邮政编码
	PostCode string `json:"postCode" gorm:"column:post_code"`
	// Province 省份/直辖市
	Province string `json:"province" gorm:"column:province"`
	// City 城市
	City string `json:"city" gorm:"column:city"`
	// Region 区
	Region string `json:"region" gorm:"column:region"`
	// DetailAddress 详细地址(街道)
	DetailAddress string `json:"detailAddress" gorm:"column:detail_address"`
}

func (address *UmsMemberReceiveAddress) GetAddressByid(id int64) {
	global.Db.Model(&UmsMemberReceiveAddress{}).Where("id = ?", id).Find(&address)
}
func (UmsMemberReceiveAddress) TableName() string {
	return "ums_member_receive_address"
}
