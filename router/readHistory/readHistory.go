package readHistory

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/readHistory"
)

type ReadHistoryRouter struct{}

func (c *ReadHistoryRouter) InitReadHistoryRouter(Router *gin.RouterGroup) {
	readHistoryRouter := Router.Group("member").Use()
	{
		readHistoryController := new(readHistory.ReadHistoryController)
		readHistoryRouter.POST("/readHistory/create", readHistoryController.Create)
		readHistoryRouter.POST("/readHistory/clear", readHistoryController.Clear)
		readHistoryRouter.POST("/readHistory/delete", readHistoryController.Delete)
		readHistoryRouter.GET("/readHistory/list", readHistoryController.List)
	}
}
