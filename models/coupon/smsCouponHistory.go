package coupon

import (
	"gomall/global"
	"time"
)

type SmsCouponHistory struct {
	ID             int64      `json:"id" gorm:"column:id"`
	CouponID       int64      `json:"couponId" gorm:"column:coupon_id"`
	MemberID       int64      `json:"memberId" gorm:"column:member_id"`
	CouponCode     string     `json:"couponCode" gorm:"column:coupon_code"`
	MemberNickname string     `json:"memberNickname" gorm:"column:member_nickname"`
	GetType        int        `json:"getType" gorm:"column:get_type"`
	CreateTime     time.Time  `json:"createTime" gorm:"column:create_time"`
	UseStatus      int        `json:"useStatus" gorm:"column:use_status"`
	UseTime        *time.Time `json:"useTime" gorm:"column:use_time"`
	OrderID        int64      `json:"orderId" gorm:"column:order_id"`
	OrderSn        string     `json:"orderSn" gorm:"column:order_sn"`
}

type SmsCouponHistoryList []SmsCouponHistory

// 这种赋值是不需要写指针的。find方法期望接收到一个切片类型的参数
func (hist *SmsCouponHistoryList) GetByCouponIdAndMemberId(couponId int64, memberId int64) (err error) {
	return global.Db.Model(&SmsCouponHistory{}).
		Where("coupon_id =? and member_id =?", couponId, memberId).Find(hist).Error
}

func (hist *SmsCouponHistory) Create() error {
	return global.Db.Create(hist).Error
}

func (SmsCouponHistory) TableName() string {
	return "sms_coupon_history"
}
