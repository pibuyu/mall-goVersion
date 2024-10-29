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
		productCollectionRouter.POST("/add", productCollectionController.Add)
	}
}
