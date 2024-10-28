package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type AddressRouter struct {
}

func (s *AddressRouter) InitAddressRouter(Router *gin.RouterGroup) {
	addressRouter := Router.Group("member").Use()
	{
		addressController := new(users.AddressController)
		addressRouter.GET("/address", addressController.GetAddressById)  //获取收获地址详情
		addressRouter.POST("/address/add", addressController.AddAddress) //添加收获地址
		addressRouter.POST("/address/delete", addressController.DeleteAddress)
		addressRouter.GET("/address/list", addressController.List)
		addressRouter.POST("/address/update", addressController.UpdateAddress)
	}
}
