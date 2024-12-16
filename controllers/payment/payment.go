package payment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"gomall/consts"
	controller "gomall/controllers"
	"gomall/global"
	"gomall/models/order"
	"gomall/utils/uniqueid"
	"net/http"
	"strconv"
)

type PaymentController struct {
	controller.BaseControllers
}

func (c *PaymentController) WebPay(ctx *gin.Context) {
	var err error
	outTradeNo := xid.Next()                //生成支付流水单号
	orderSn := ctx.Query("outTradeNo")      //订单号,可以拿订单号和ctx获取的memberId倒查订单的其他信息
	subject := ctx.Query("subject")         //订单注释字段
	totalAmount := ctx.Query("totalAmount") //应付金额

	//1.生成预支付单
	//todo:支付状态发生变化时需要考虑将这个hash里订单号与memberId的映射删除,避免空间无限膨胀
	memberIdFromCtx, err := global.RedisDb.HGet(consts.Order2MemberIdMap, orderSn).Result()
	memberIdInt64, _ := strconv.ParseInt(memberIdFromCtx, 10, 64)
	if err != nil {
		global.Logger.Errorf("预支付时，查询memberId failed:%v", err)
		c.Response(ctx, "获取用户信息失败", nil, err)
		return
	}
	orderInfo := &order.OmsOrder{}
	if err := global.Db.Model(&order.OmsOrder{}).Where("member_id = ? and order_sn = ?", memberIdInt64, orderSn).Find(orderInfo).Error; err != nil {
		global.Logger.Errorf("预支付时，查询订单详情failed : %v", err)
		c.Response(ctx, "查询订单详情失败", nil, err)
		return
	}
	payment := &order.OmsPayment{}
	payment.OrderId = orderInfo.ID
	payment.PaymentSn = uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	payment.UserId = memberIdInt64
	//pay_status 和 pay_type不用设置，创建预支付单时用默认的 未支付 和 支付宝方式
	payment.OrderSn = orderSn
	payment.Subject = subject
	payment.OutTradeNo = strconv.FormatInt(outTradeNo, 10)
	float64Price, _ := strconv.ParseFloat(totalAmount, 64)
	payment.TotalAmount = float64Price
	if err = payment.Insert(); err != nil {
		global.Logger.Errorf("预支付时，插入预支付单失败:%v", err)
		c.Response(ctx, "预支付时，插入预支付单失败", nil, err)
		return
	}

	//2.构造并发起支付请求
	//如果没有成功支付，退回来的话会重新生成一个本地的payment_sn，也会生成一个新的out_trad_no.
	var p = alipay.TradePagePay{}
	//这是ngrok配置的9090端口的内网穿透地址，每次重启ngrok都要修改这里的回调地址
	p.NotifyURL = "http://101.126.144.39:9090/alipay/callback"
	//支付完成之后跳转页面，直接填前端的订单页即可
	p.ReturnURL = "http://localhost:8060/#/pages/order/order?state=0"
	p.Subject = subject
	p.OutTradeNo = strconv.FormatInt(outTradeNo, 10)
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" //网页支付,这个字段是固定的，目前支付沙箱仅支持这一个模式

	url, _ := global.AliPay.TradePagePay(p)
	global.Logger.Infof("订单:%s尝试发起支付请求", p.OutTradeNo)
	//跳转到支付宝支付页面
	ctx.Redirect(http.StatusTemporaryRedirect, url.String())
}

func (c *PaymentController) AliPayCallback(ctx *gin.Context) {
	//用密钥鉴别一下这个回调是不是支付宝发过来的
	if err := ctx.Request.ParseForm(); err != nil {
		global.Logger.Errorf("解析支付回调参数failed:%v", err)
		c.Response(ctx, "解析支付回调参数失败", nil, err)
		return
	}

	// 获取通知参数
	outTradeNo := ctx.Request.Form.Get("out_trade_no")
	// 验证签名，检查请求是否安全
	if err := global.AliPay.VerifySign(ctx.Request.Form); err != nil {
		global.Logger.Errorf("验证支付回调签名未通过：%v", err)
		c.Response(ctx, "验证支付回调签名未通过", nil, err)
		return
	}
	// 主动查询支付状态
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	resp, err := global.AliPay.TradeQuery(ctx, p)
	if err != nil {
		global.Logger.Errorf("查询订单支付状态失败：%v", err)
		c.Response(ctx, "查询订单支付状态失败", nil, err)
		return
	}

	// 支付失败
	if resp.TradeStatus != "TRADE_SUCCESS" {
		global.Logger.Errorf("支付流水单：%s 支付失败", outTradeNo)
		c.Response(ctx, "订单支付失败", nil, errors.New("订单支付失败"))
		return
	}

	// 根据outTradeNo查询本地的支付记录，比对支付的金额是否一致
	payment := &order.OmsPayment{}
	if err := global.Db.Model(&order.OmsPayment{}).Where("out_trade_no = ?", outTradeNo).Find(&payment).Error; err != nil {
		global.Logger.Errorf("根据payment_sn查询订单详情failed : %v", err)
		c.Response(ctx, "查询订单详情失败", nil, err)
		return
	}
	payAmount, _ := strconv.ParseFloat(resp.TotalAmount, 10)
	if payAmount != payment.TotalAmount {
		global.Logger.Errorf("支付流水%s 的实际支付金额= %v ,应付金额为= %v ", outTradeNo, payAmount, payment.TotalAmount)
		c.Response(ctx, "实付金额与应付金额不符", nil, err)
		return
	}

	//支付成功,修改订单状态为已支付
	if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", payment.OrderId).UpdateColumn("status", 1).Error; err != nil {
		global.Logger.Errorf("修改订单 %d 为已支付状态失败 %v", payment.OrderId, err)
	}

	//修改支付状态、实付金额
	payment.PayStatus = 1
	payment.ActualPay = payAmount
	if err = global.Db.Model(&order.OmsPayment{}).Where("out_trade_no = ?", payment.OutTradeNo).Updates(payment).Error; err != nil {
		global.Logger.Errorf("修改支付流水 %s 出错 %v", outTradeNo, err)
	}
	c.Response(ctx, "订单支付成功", nil, nil)
}

// 不要等结果，如果真的修改失败了通知用户即可
func updatePaymentStatus(outTradeNo string, orderId int64, payAmount float64, newStatus int) {
	var err error
	tx := global.Db.Begin()

	if err = tx.Table("oms_order").Debug().Where("id = ?", orderId).UpdateColumn("status", 1).Error; err != nil {
		global.Logger.Errorf("修改订单 : %d 状态failed:%v", orderId, err)
		tx.Rollback()
	}
	if err = tx.Table("oms_payment").Debug().Where("out_trade_no = ?", outTradeNo).UpdateColumn("pay_status", 1).UpdateColumn("actual_pay", payAmount).Error; err != nil {
		global.Logger.Errorf("修改订单 : %s 状态failed:%v", outTradeNo, err)
		tx.Rollback()
		//todo:应该丢进消息队列一直重试
	}
}
