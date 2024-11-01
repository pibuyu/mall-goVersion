package home

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/home"
)

type HomeRouter struct{}

func (c *HomeRouter) InitHomeRouter(Router *gin.RouterGroup) {
	homeRouter := Router.Group("home").Use()
	{
		homeController := new(home.HomeController)
		homeRouter.GET("/content", homeController.GetHomeContent)
		homeRouter.GET("/hotProductList", homeController.GetHotProductList)
		homeRouter.GET("/newProductList", homeController.GetNewProductList)
		homeRouter.GET("/recommendProductList", homeController.GetRecommendProductList)
		homeRouter.GET("/subjectList", homeController.GetSubjectList)
		homeRouter.GET("/productCateList/:parentId", homeController.GetProductCateList)
	}
}
