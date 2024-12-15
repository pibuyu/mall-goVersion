package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	controller "gomall/controllers"
	"gomall/global"
	"net/http"
)

type PaymentController struct {
	controller.BaseControllers
}

func (c *PaymentController) WebPay(ctx *gin.Context) {

	//todo:应该由前端传递过来支付流水号，订单金额参数
	outTradeNo := ctx.Query("outTradeNo")
	subject := ctx.Query("subject")
	totalAmount := ctx.Query("totalAmount")
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http:101.126.144.39:9090/alipay/callback"
	p.ReturnURL = "https://6001-2001-250-3000-2b70-6459-b39e-d303-267e.ngrok-free.app/#/pages/order/order?state=0"
	p.Subject = subject
	p.OutTradeNo = outTradeNo
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" //网页支付,这个字段是固定的，目前支付沙箱仅支持这一个模式

	url, _ := global.AliPay.TradePagePay(p)
	var payUrl = url.String()
	global.Logger.Infof("即将跳转到支付宝页面:%s", payUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, url.String())
}

func (c *PaymentController) Callback(ctx *gin.Context) {
	global.Logger.Infoln("收到了支付宝的支付回调")
}
