package readHistory

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/readHistory"
)

type ReadHistoryRouter struct{}

func (c *ReadHistoryRouter) InitReadHistoryRouter(Router *gin.RouterGroup) {
	readHistoryRouter := Router.Group("/member").Use()
	{
		//和浏览记录相关的操作需要放在mongodb，而不是mysql中
		readHistoryController := new(readHistory.ReadHistoryController)
		readHistoryRouter.GET("/readHistory/list", readHistoryController.List)
		readHistoryRouter.POST("/readHistory/clear", readHistoryController.Clear)
		readHistoryRouter.POST("/readHistory/create", readHistoryController.Create)
		readHistoryRouter.POST("/readHistory/delete", readHistoryController.Delete)
	}
}
