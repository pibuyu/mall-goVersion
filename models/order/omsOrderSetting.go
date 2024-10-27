package order

import (
	"errors"
	"gomall/global"
)

type OmsOrderSetting struct {
	// ID 记录ID
	ID int64 `json:"id" gorm:"column:id"`
	// FlashOrderOvertime 秒杀订单超时关闭时间（分）
	FlashOrderOvertime int `json:"flashOrderOvertime" gorm:"column:flash_order_overtime"`
	// NormalOrderOvertime 正常订单超时时间（分）
	NormalOrderOvertime int `json:"normalOrderOvertime" gorm:"column:normal_order_overtime"`
	// ConfirmOvertime 发货后自动确认收货时间（天）
	ConfirmOvertime int `json:"confirmOvertime" gorm:"column:confirm_overtime"`
	// FinishOvertime 自动完成交易时间，不能申请售后（天）
	FinishOvertime int `json:"finishOvertime" gorm:"column:finish_overtime"`
	// CommentOvertime 订单完成后自动好评时间（天）
	CommentOvertime int `json:"commentOvertime" gorm:"column:comment_overtime"`
}
type OmsOrderSettingList []OmsOrderSetting

func (settingList *OmsOrderSettingList) GetAll() (err error) {
	if err := global.Db.Find(&settingList).Error; err != nil {
		return errors.New("查询全部orderSettings出错:" + err.Error())
	}
	return nil
}

func (OmsOrderSetting) TableName() string {
	return "oms_order_setting"
}
