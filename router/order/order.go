package order

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gomall/controllers/order"
	"gomall/utils/jwt"
	"net/http"
	"sync"
)

type OrderRouter struct{}

var userLimiters sync.Map

func (c *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("order").Use()
	{
		orderController := new(order.OrderController)
		orderRouter.GET("/detail/:orderId", orderController.Detail)                     //根据id获取订单详情
		orderRouter.POST("/cancelOrder", orderController.CancelOrder)                   //取消单个超时订单
		orderRouter.POST("/confirmReceiveOrder", orderController.ConfirmReceiveOrder)   //确认收货
		orderRouter.POST("/deleteOrder", orderController.DeleteOrder)                   //删除订单
		orderRouter.POST("/generateConfirmOrder", orderController.GenerateConfirmOrder) //根据购物车信息生成确认单
		//添加令牌桶限流
		orderRouter.POST("/generateOrder", c.LimitGenerateOrder(orderController.GenerateOrder)) //根据购物车信息生成订单
		orderRouter.GET("/list", orderController.List)                                          //根据id获取订单详情
		orderRouter.POST("/paySuccess", orderController.PaySuccess)                             //支付成功的回调
		orderRouter.POST("/cancelUserOrder", orderController.CancelUserOrder)                   //用户取消订单
		orderRouter.POST("/cancelTimeOutOrder", orderController.CancelTimeOutOrder)             //自动取消超时订单
	}

	//单独的一条退货申请的路由
	orderController := new(order.OrderController)
	Router.POST("/returnApply/create", orderController.CreateReturnApply)
}

func (c *OrderRouter) LimitGenerateOrder(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, _ := jwt.GetMemberIdFromCtx(c)
		limiter, _ := userLimiters.LoadOrStore(memberId, rate.NewLimiter(1, 3)) //每秒每个用户最多访问1次
		if !limiter.(*rate.Limiter).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "您的操作过于频繁，请稍后再试"})
			return
		}
		handlerFunc(c) // 继续处理请求
	}
}
