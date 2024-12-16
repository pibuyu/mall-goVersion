package payment

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/payment"
)

type PaymentRouter struct {
}

func (c *PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {
	paymentRouter := Router.Group("/alipay").Use()
	{
		paymentController := new(payment.PaymentController)
		paymentRouter.GET("/webPay", paymentController.WebPay)
		paymentRouter.POST("/callback", paymentController.AliPayCallback)
	}
}
