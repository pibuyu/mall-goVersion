package order

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/order"
	"gomall/logic/order"
	"gomall/utils/jwt"
	"gomall/utils/response"
	"strconv"
)

type OrderController struct {
	controller.BaseControllers
}

func (c *OrderController) Detail(ctx *gin.Context) {
	orderId, _ := strconv.ParseInt(ctx.Param("orderId"), 10, 64)
	result, err := order.Detail(orderId)
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
	if err := order.CancelOrder(rec.OrderId); err != nil {
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
	rec.OrderId, _ = strconv.ParseInt(ctx.PostForm("orderId"), 10, 64)
	if err := order.DeleteOrder(ctx, &rec); err != nil {
		global.Logger.Error("删除订单出错:%v", err)
		c.Response(ctx, "删除订单失败", nil, err)
		return
	}
	c.Response(ctx, "删除订单成功", nil, nil)
}

// 根据购物车信息生成确认单
func (c *OrderController) GenerateConfirmOrder(ctx *gin.Context) {
	var ids []int64
	if err := ctx.ShouldBind(&ids); err != nil {
		global.Logger.Errorf("GenerateConfirmOrder请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	confirmOrder, err := order.GenerateConfirmOrder(ids, ctx)
	if err != nil {
		global.Logger.Error("根据购物车信息生成确认单出错:%v", err)
		c.Response(ctx, "根据购物车信息生成确认单失败", nil, err)
		return
	}
	c.Response(ctx, "根据购物车信息生成确认单成功", confirmOrder, nil)

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
	if err := ctx.ShouldBind(&rec); err != nil {
		global.Logger.Errorf("List请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	result, err := order.List(&rec, memberId)
	if err != nil {
		c.Response(ctx, "按订单状态分页获取订单列表失败", nil, err)
	}
	pageResult := response.ResetPage(result, int64(len(result)), rec.PageNum, rec.PageSize)
	c.Response(ctx, "按订单状态分页获取订单列表成功", pageResult, nil)
}

// PaySuccess 用户支付成功的回调
func (c *OrderController) PaySuccess(ctx *gin.Context) {
	var rec receive.PaySuccessReqStruct
	rec.OrderId, _ = strconv.ParseInt(ctx.PostForm("orderId"), 10, 64)
	rec.PayType, _ = strconv.Atoi(ctx.PostForm("payType"))

	count, err := order.PaySuccess(&rec)
	if err != nil {
		c.Response(ctx, "支付失败", nil, err)
		return
	}
	c.Response(ctx, "支付成功", count, nil)
}

// CancelUserOrder 用户取消订单
func (c *OrderController) CancelUserOrder(ctx *gin.Context) {
	orderId, _ := strconv.ParseInt(ctx.PostForm("orderId"), 10, 64)

	if err := order.CancelOrder(orderId); err != nil {
		global.Logger.Error("用户主动取消订单出错:%v", err)
		c.Response(ctx, "用户主动取消订单失败", nil, err)
		return
	}
	c.Response(ctx, "用户主动取消订单成功", nil, nil)
}

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
