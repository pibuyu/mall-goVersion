package router

import (
	"github.com/gin-gonic/gin"
	middlewares "gomall/middleware"
	cartRouter "gomall/router/cart"
	homeRouter "gomall/router/home"
	orderRouter "gomall/router/order"
	productRouter "gomall/router/product"
	readHistoryRouter "gomall/router/readHistory"
	usersRouter "gomall/router/users"
)

type RouterGroup struct {
	Users         usersRouter.RouterGroup
	Home          homeRouter.RouterGroup
	Cart          cartRouter.RouterGroup
	ReadHistory   readHistoryRouter.RouterGroup
	OrderRouter   orderRouter.RouterGroup
	ProductRouter productRouter.RouterGroup
}

var RoutersGroup = new(RouterGroup)

func InitRouter() {
	router := gin.Default()
	//跨域中间件，允许所有源
	router.Use(middlewares.Cors())
	PrivateGroup := router.Group("")
	PrivateGroup.Use()
	{
		//静态资源
		router.Static("/assets", "./assets")
		//初始化各个路由器组
		RoutersGroup.Users.LoginRouter.InitLoginRouter(PrivateGroup)
		RoutersGroup.Users.AddressRouter.InitAddressRouter(PrivateGroup)
		RoutersGroup.Users.CouponRouter.InitCouponRouter(PrivateGroup)
		RoutersGroup.Users.ProductCollectionRouter.InitProductCollectionRouter(PrivateGroup)
		RoutersGroup.Users.BrandAttentionRouter.InitBrandAttentionRouter(PrivateGroup)
		RoutersGroup.Users.BrandRouter.InitBrandRouter(PrivateGroup)
		RoutersGroup.Home.HomeRouter.InitHomeRouter(PrivateGroup)
		RoutersGroup.Cart.InitCartRouter(PrivateGroup)
		RoutersGroup.ReadHistory.InitReadHistoryRouter(PrivateGroup)
		RoutersGroup.OrderRouter.InitOrderRouter(PrivateGroup)
		RoutersGroup.ProductRouter.InitProductRouter(PrivateGroup)
	}

	if err := router.Run(":9090"); err != nil {
		return
	}
}
