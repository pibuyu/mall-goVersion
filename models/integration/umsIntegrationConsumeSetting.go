package integration

import "gomall/global"

// 积分消费设置表
type UmsIntegrationConsumeSetting struct {
	// ID 记录ID
	ID int64 `json:"id" gorm:"column:id"`
	// DeductionPerAmount 每一元需要抵扣的积分数量
	DeductionPerAmount int `json:"deductionPerAmount" gorm:"column:deduction_per_amount"`
	// MaxPercentPerOrder 每笔订单最高抵用百分比
	MaxPercentPerOrder int `json:"maxPercentPerOrder" gorm:"column:max_percent_per_order"`
	// UseUnit 每次使用积分最小单位100
	UseUnit int `json:"useUnit" gorm:"column:use_unit"`
	// CouponStatus 是否可以和优惠券同用；0 -> 不可以；1 -> 可以
	CouponStatus int `json:"couponStatus" gorm:"column:coupon_status"`
}

func (UmsIntegrationConsumeSetting) TableName() string {
	return "ums_integration_consume_setting"
}

func (setting *UmsIntegrationConsumeSetting) GetById(id int64) {
	global.Db.Model(&UmsIntegrationConsumeSetting{}).Where("id = ?", id).Find(&setting)
}
