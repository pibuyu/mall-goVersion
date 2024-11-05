package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type CouponRouter struct {
}

func (s *CouponRouter) InitCouponRouter(Router *gin.RouterGroup) {
	couponRouter := Router.Group("/member/coupon").Use()
	{
		couponController := new(users.CouponController)
		couponRouter.POST("/add/:couponId", couponController.Add)
		couponRouter.GET("/list", couponController.List)
		couponRouter.GET("/list/cart", couponController.ListCart)                     //这个接口的正确性暂时无从验证
		couponRouter.GET("/listByProduct/:productId", couponController.ListByProduct) //这个接口的正确性得到了部分验证，但测试不充分
		couponRouter.GET("/listHistory", couponController.ListHistory)
	}
}
