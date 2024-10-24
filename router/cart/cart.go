package cart

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/cart"
	middlewares "gomall/middleware"
)

type CartRouter struct{}

func (c *CartRouter) InitCartRouter(Router *gin.RouterGroup) {
	cartRouter := Router.Group("cart").Use(middlewares.VerificationToken())
	{
		cartController := new(cart.CartController)
		cartRouter.POST("/add", cartController.AddCartItem)
		cartRouter.POST("/clear", cartController.Clear)
		cartRouter.POST("/delete", cartController.DeleteByIds)
		cartRouter.GET("/getProduct", cartController.GetProductById)
		cartRouter.GET("/list", cartController.List)
		//todo:下面这个接口比较复杂，正确性有待验证。暂时没找到原界面在哪里调用的，不知道构造什么请求体数据可以有效验证
		cartRouter.GET("/list/promotion", cartController.CartListPromotion)
		cartRouter.POST("/update/attr", cartController.UpdateAttr)
		cartRouter.POST("/update/quantity", cartController.UpdateQuantity)
	}
}
