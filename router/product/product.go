package product

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/product"
)

type ProductRouter struct{}

func (c *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	productRouter := Router.Group("product").Use()
	{
		productController := new(product.ProductController)
		productRouter.GET("/detail", productController.Detail)
		productRouter.GET("/categoryTreeList", productController.CategoryTreeList)
		productRouter.GET("/search", productController.Search)
	}
}
