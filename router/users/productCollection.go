package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type ProductCollectionRouter struct {
}

func (r *ProductCollectionRouter) InitProductCollectionRouter(Router *gin.RouterGroup) {
	productCollectionRouter := Router.Group("/member/productCollection").Use()
	{
		productCollectionController := new(users.ProductCollectionController)
		//todo:这里也全都是操作mongodb
		productCollectionRouter.POST("/add", productCollectionController.Add)
		productCollectionRouter.POST("/clear", productCollectionController.Clear)
		productCollectionRouter.POST("/delete", productCollectionController.Delete)
		productCollectionRouter.GET("/detail", productCollectionController.Detail)
		productCollectionRouter.GET("/list", productCollectionController.List)
	}
}
