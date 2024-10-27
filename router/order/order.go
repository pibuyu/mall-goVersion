package order

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/order"
)

type OrderRouter struct{}

func (c *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("order").Use()
	{
		orderController := new(order.OrderController)
		orderRouter.GET("/detail", orderController.Detail) //根据id获取订单详情
		//这个接口暂时无从验证对错
		orderRouter.POST("/cancelOrder", orderController.CancelOrder)
		orderRouter.POST("/confirmReceiveOrder", orderController.ConfirmReceiveOrder)   //确认收货
		orderRouter.POST("/deleteOrder", orderController.DeleteOrder)                   //删除订单
		orderRouter.POST("/generateConfirmOrder", orderController.GenerateConfirmOrder) //根据购物车信息生成确认单
		orderRouter.POST("/generateOrder", orderController.GenerateOrder)               //根据购物车信息生成确认单
	}
}