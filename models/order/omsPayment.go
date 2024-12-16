package order

import "gomall/global"

type OmsPayment struct {
	Id          int64   `json:"id" gorm:"id"`
	OrderId     int64   `json:"order_id" gorm:"order_id"`
	UserId      int64   `json:"user_id" gorm:"user_id"`
	PaymentSn   string  `json:"payment_sn" gorm:"payment_sn"`     //支付流水单号
	ActualPay   float64 `json:"actual_pay" gorm:"actual_pay"`     //实际支付金额
	OrderSn     string  `json:"order_sn" gorm:"order_sn"`         //订单sn码
	PayType     string  `json:"pay_type" gorm:"pay_type"`         //支付方式，现在默认为1，表示支付宝支付
	PayStatus   int     `json:"pay_status" gorm:"pay_status"`     //支付状态：0-未支付；1-正在支付中；2-支付取消；3-支付成功
	OutTradeNo  string  `json:"out_trade_no" gorm:"out_trade_no"` //前端传递的，生成的订单号
	Subject     string  `json:"subject" gorm:"subject"`           //前端传递的，订单注释字段
	TotalAmount float64 `json:"total_amount" gorm:"total_amount"` //前端传递的，应付金额
}

func (OmsPayment) TableName() string {
	return "oms_payment"
}

func (payment *OmsPayment) Insert() error {
	return global.Db.Model(&OmsPayment{}).Create(payment).Error
}
