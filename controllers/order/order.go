package order

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/order"
	"gomall/logic/order"
)

type OrderController struct {
	controller.BaseControllers
}

func (c *OrderController) Detail(ctx *gin.Context) {
	var rec receive.DetailReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("Detail请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	result, err := order.Detail(&rec)
	if err != nil {
		c.Response(ctx, "获取订单详情失败", nil, err)
	}
	c.Response(ctx, "获取订单详情成功", result, nil)
}

// 全程不给用户报错，体验不好，后台处理就行
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	var rec receive.CancelOrderReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("Detail请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	if err := order.CancelOrder(&rec); err != nil {
		global.Logger.Error("取消订单出错:%v", err)
	}

	c.Response(ctx, "取消订单成功", nil, nil)
}

func (c *OrderController) ConfirmReceiveOrder(ctx *gin.Context) {
	var rec receive.ConfirmReceiveOrderReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("ConfirmReceiveOrder请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	if err := order.ConfirmReceiveOrder(ctx, &rec); err != nil {
		global.Logger.Error("确认收货出错:%v", err)
		c.Response(ctx, "确认收货订单失败", nil, err)
		return
	}
	c.Response(ctx, "确认收货成功", nil, nil)
}

func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	var rec receive.DeleteOrderReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("DeleteOrder请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	if err := order.DeleteOrder(ctx, &rec); err != nil {
		global.Logger.Error("删除订单出错:%v", err)
		c.Response(ctx, "删除订单失败", nil, err)
		return
	}
	c.Response(ctx, "删除订单成功", nil, nil)
}

// 根据购物车信息生成确认单
func (c *OrderController) GenerateConfirmOrder(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GenerateConfirmOrderReqStruct)); err == nil {
		confirmOrder, err := order.GenerateConfirmOrder(rec.CartIds, ctx)
		if err != nil {
			global.Logger.Error("根据购物车信息生成确认单出错:%v", err)
			c.Response(ctx, "根据购物车信息生成确认单失败", nil, err)
			return
		}
		c.Response(ctx, "根据购物车信息生成确认单成功", confirmOrder, nil)
	}
}

func (c *OrderController) GenerateOrder(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.GenerateOrderReqStruct)); err == nil {
		order, err := order.GenerateOrder(rec, ctx)
		if err != nil {
			global.Logger.Error("根据购物车信息生成订单出错:%v", err)
			c.Response(ctx, "根据购物车信息生成订单失败", nil, err)
			return
		}
		c.Response(ctx, "根据购物车信息生成订单成功", order, nil)
	}
}
