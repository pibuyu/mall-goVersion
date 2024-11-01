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
		addressRouter.GET("/address/:addressId", addressController.GetAddressById) //获取收获地址详情
		addressRouter.POST("/address/add", addressController.AddAddress)           //添加收获地址
		addressRouter.POST("/address/delete/:addressId", addressController.DeleteAddress)
		addressRouter.GET("/address/list", addressController.List)
		addressRouter.POST("/address/update/:addressId", addressController.UpdateAddress) //这个接口更新成功前端没反应的。。还以为更新失败了
	}
}
