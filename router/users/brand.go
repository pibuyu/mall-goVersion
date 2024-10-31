package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type BrandRouter struct {
}

func (s *BrandRouter) InitBrandRouter(Router *gin.RouterGroup) {
	brandRouter := Router.Group("/brand").Use()
	{
		brandController := new(users.BrandController)
		brandRouter.GET("/detail", brandController.Detail)
		brandRouter.GET("/recommendList", brandController.RecommendList)
		brandRouter.GET("/productList", brandController.ProductList)
	}
}
