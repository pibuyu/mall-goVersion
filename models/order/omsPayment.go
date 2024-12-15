package order

import "gomall/global"

type OmsPayment struct {
	Id        int64   `json:"id" gorm:"id"`
	UserId    int64   `json:"user_id" gorm:"user_id"`
	OrderSn   string  `json:"order_sn" gorm:"order_sn"`
	PayType   string  `json:"pay_type" gorm:"pay_type"`
	PayStatus int     `json:"pay_status" gorm:"pay_status"`
	OrderId   string  `json:"order_id" gorm:"order_id"`
	Price     float32 `json:"price" gorm:"price"`
}

func (OmsPayment) TableName() string {
	return "oms_payment"
}

func (payment *OmsPayment) Insert() error {
	return global.Db.Model(&OmsPayment{}).Create(payment).Error
}
