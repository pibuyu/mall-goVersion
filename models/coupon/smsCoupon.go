package coupon

import (
	"gomall/global"
	"time"
)

// SmsCoupon 结构体定义，用于表示优惠券相关信息
type SmsCoupon struct {
	// ID 优惠券记录ID
	ID int64 `json:"id" gorm:"id"`
	// Type 优惠券类型；0->全场赠券；1->会员赠券；2->购物赠券；3->注册赠券
	Type int `json:"type" gorm:"type"`
	// Name 优惠券名称
	Name string `json:"name" gorm:"name"`
	// Platform 使用平台：0->全部；1->移动；2->PC
	Platform int `json:"platform" gorm:"column:platform"`
	// Count 数量
	Count int `json:"count" gorm:"column:count"`
	// Amount 金额
	Amount float32 `json:"amount" gorm:"amount"`
	// PerLimit 每人限领张数
	PerLimit int `json:"perLimit" gorm:"per_limit"`
	// MinPoint 使用门槛；0表示无门槛
	MinPoint float32 `json:"minPoint" gorm:"column:min_point"`
	// StartTime 开始时间
	StartTime time.Time `json:"startTime" gorm:"column:start_time"`
	// EndTime 结束时间
	EndTime time.Time `json:"endTime" gorm:"column:end_time"`
	// UseType 使用类型：0->全场通用；1->指定分类；2->指定商品
	UseType int `json:"useType" gorm:"column:use_type"`
	// Note 备注
	Note string `json:"note" gorm:"column:note"`
	// PublishCount 发行数量
	PublishCount int `json:"publishCount" gorm:"column:publish_count"`
	// UseCount 已使用数量
	UseCount int `json:"useCount" gorm:"column:use_count"`
	// ReceiveCount 领取数量
	ReceiveCount int `json:"receiveCount" gorm:"column:receive_count"`
	// EnableTime 可以领取的日期
	EnableTime time.Time `json:"enableTime" gorm:"column:enable_time"`
	// Code 优惠码
	Code string `json:"code" gorm:"column:code"`
	// MemberLevel 可领取的会员类型：0->无限时
	MemberLevel int `json:"memberLevel" gorm:"column:member_level"`
}

type SmsCouponList []SmsCoupon

func (list *SmsCouponList) GetListByIds(couponIds []int64) error {
	query := global.Db.Model(&SmsCoupon{}).
		Where("end_time > ? and start_time < ? and use_type = ?", time.Now(), time.Now(), 0)

	if len(couponIds) > 0 {
		query = query.Or("end_time > ? AND start_time < ? AND use_type != ? AND id IN ?",
			time.Now(), time.Now(), 0, couponIds)
	}
	return query.Find(list).Error
}

func (c *SmsCoupon) Update() error {
	return global.Db.Model(&SmsCoupon{}).Where("id", c.ID).Updates(c).Error
}

func (c *SmsCoupon) GetCouponById(couponId int64) (err error) {
	return global.Db.Model(&SmsCoupon{}).Where("id = ?", couponId).Find(c).Error
}

func (SmsCoupon) TableName() string {
	return "sms_coupon"
}
