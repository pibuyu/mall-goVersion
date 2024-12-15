package es

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/es"
)

type EsProductRouter struct{}

func (c *EsProductRouter) InitEsProductRouter(Router *gin.RouterGroup) {
	esProductRouter := Router.Group("esProduct").Use()
	{
		esProductController := new(es.EsProductController)
		esProductRouter.POST("/importAll", esProductController.ImportAll)
		esProductRouter.POST("/delete/:id", esProductController.Delete)
		esProductRouter.POST("/create/:id", esProductController.Create)
		esProductRouter.GET("/search/simple", esProductController.SimpleSearch) //todo:这个地方查不到结果
	}
}
