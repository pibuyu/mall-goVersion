package users

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/address"
	addressLogic "gomall/logic/address"
	"gomall/utils/jwt"
)

type AddressController struct {
	controller.BaseControllers
}

// 获取收货地址详情
func (c *AddressController) GetAddressById(ctx *gin.Context) {
	var rec receive.GetAddressByIdReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("根据id获取地址时，绑定参数错误：%v", err)
		c.Response(ctx, "根据id获取地址时，绑定参数错误", nil, err)
		return
	}

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	result, err := addressLogic.GetAddressById(rec.Id, memberId)
	if err != nil {
		c.Response(ctx, "根据id获取地址时，查表错误", nil, err)
		return
	}
	c.Response(ctx, "根据id获取地址成功", result, nil)
}

func (c *AddressController) AddAddress(ctx *gin.Context) {
	var rec receive.AddAddressReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("添加地址时，绑定参数错误：%v", err)
		c.Response(ctx, "添加地址时，绑定参数错误", nil, err)
		return

	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	rec.MemberId = memberId

	if err := addressLogic.AddAddress(&rec); err != nil {
		c.Response(ctx, "添加地址时，插入数据库错误", 0, err)
		return
	}

	c.Response(ctx, "添加地址成功", 1, nil)
}

func (c *AddressController) DeleteAddress(ctx *gin.Context) {
	var rec receive.DeleteAddressReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("删除地址时，绑定参数错误：%v", err)
		c.Response(ctx, "删除地址时，绑定参数错误", nil, err)
		return
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	if err := addressLogic.DeleteAddress(rec.Id, memberId); err != nil {
		c.Response(ctx, "删除地址时，删除数据库表项出错", 0, err)
		return
	}

	c.Response(ctx, "删除地址成功", 1, nil)
}
func (c *AddressController) List(ctx *gin.Context) {
	result, err := addressLogic.List()
	if err != nil {
		c.Response(ctx, "获取地址列表failed:", nil, err)
	}

	c.Response(ctx, "获取地址列表成功", result, nil)
}

func (c *AddressController) UpdateAddress(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UpdateAddressReqStruct)); err == nil {
		memberId, _ := jwt.GetMemberIdFromCtx(ctx)
		if err := addressLogic.UpdateAddress(rec, memberId); err != nil {
			global.Logger.Errorf("更新收货地址failed:" + err.Error())
			c.Response(ctx, "更新收货地址出错", 0, err)
			return
		}
		c.Response(ctx, "更新收货地址成功", 1, nil)
	}
}
