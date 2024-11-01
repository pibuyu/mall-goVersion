package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type BrandAttentionRouter struct {
}

func (s *BrandAttentionRouter) InitBrandAttentionRouter(Router *gin.RouterGroup) {
	brandAttentionRouter := Router.Group("/member/attention").Use()
	{
		brandAttentionController := new(users.BrandAttentionController)
		brandAttentionRouter.GET("/list", brandAttentionController.List)
		brandAttentionRouter.POST("/add", brandAttentionController.Add)
		brandAttentionRouter.POST("/clear", brandAttentionController.Clear)
		brandAttentionRouter.POST("/delete", brandAttentionController.Delete)
		brandAttentionRouter.GET("/detail", brandAttentionController.Detail)
	}
}
