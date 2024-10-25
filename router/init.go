package router

import (
	"github.com/gin-gonic/gin"
	middlewares "gomall/middleware"
	cartRouter "gomall/router/cart"
	homeRouter "gomall/router/home"
	readHistoryRouter "gomall/router/readHistory"
	usersRouter "gomall/router/users"
)

type RouterGroup struct {
	Users       usersRouter.RouterGroup
	Home        homeRouter.RouterGroup
	Cart        cartRouter.RouterGroup
	ReadHistory readHistoryRouter.RouterGroup
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
		RoutersGroup.Home.HomeRouter.InitHomeRouter(PrivateGroup)
		RoutersGroup.Cart.InitCartRouter(PrivateGroup)
		RoutersGroup.ReadHistory.InitReadHistoryRouter(PrivateGroup)
	}

	if err := router.Run(":9090"); err != nil {
		return
	}
}
