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
		brandAttentionRouter.POST("/add", brandAttentionController.Add) //添加收获地址
	}
}
