package order

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/order"
	"gomall/logic/order"
	"gomall/utils/jwt"
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
		oneOrder, err := order.GenerateOrder(rec, ctx)
		if err != nil {
			global.Logger.Error("根据购物车信息生成订单出错:%v", err)
			c.Response(ctx, "根据购物车信息生成订单失败", nil, err)
			return
		}
		c.Response(ctx, "根据购物车信息生成订单成功", oneOrder, nil)
	}
}

// List 按订单状态分页获取订单列表
func (c *OrderController) List(ctx *gin.Context) {
	var rec receive.ListReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("List请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	result, err := order.List(&rec, memberId)
	if err != nil {
		c.Response(ctx, "按订单状态分页获取订单列表失败", nil, err)
	}
	c.Response(ctx, "按订单状态分页获取订单列表成功", result, nil)
}

// PaySuccess 用户支付成功的回调
func (c *OrderController) PaySuccess(ctx *gin.Context) {
	var rec receive.PaySuccessReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("PaySuccess请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}

	count, err := order.PaySuccess(&rec)
	if err != nil {
		c.Response(ctx, "支付失败", nil, err)
		return
	}
	c.Response(ctx, "支付成功", count, nil)
}

// CancelUserOrder 用户取消订单
func (c *OrderController) CancelUserOrder(ctx *gin.Context) {
	var rec receive.CancelUserOrderReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("CancelUserOrder请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}

	newRec := &receive.CancelOrderReqStruct{
		OrderId: rec.OrderId,
	}
	if err := order.CancelOrder(newRec); err != nil {
		global.Logger.Error("用户主动取消订单出错:%v", err)
		c.Response(ctx, "用户主动取消订单失败", nil, err)
		return
	}
	c.Response(ctx, "用户主动取消订单成功", nil, nil)
}

// todo:这个方法应该由定时扫描的定时器来做.如果是手动调用，相当于手动清理那些超时订单
func (c *OrderController) CancelTimeOutOrder(ctx *gin.Context) {
	//不需要参数，直接扫描超时订单
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	count, err := order.CancelTimeOutOrder(memberId)
	if err != nil {
		global.Logger.Error("取消超时订单出错:%v", err)
		c.Response(ctx, "取消超时订单失败", nil, err)
		return
	}
	c.Response(ctx, "取消超时订单成功", count, nil)
}

func (c *OrderController) CreateReturnApply(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.CreateReturnApplyReqStruct)); err == nil {
		if err := order.CreateReturnApply(rec); err != nil {
			c.Response(ctx, "创建退货订单failed", 0, err)
		}
		c.Response(ctx, "创建退货订单成功", 1, nil)
	}
}
