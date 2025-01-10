package order

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/order"
	"gomall/logic/order"
	order2 "gomall/models/order"
	"gomall/utils/jwt"
	"gomall/utils/response"
	"gomall/utils/uniqueid"
	"os"
	"runtime/pprof"
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
	//创建pprof文件
	file, err := os.Create("generateOrderCpu.pprof")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	//开始cpu性能采样
	_ = pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

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
	//todo:go-zero-look-look里的流程是：先拿userId换openId，这里没有openId，跳过。
	//然后生成本地的支付流水号
	//然后创建微信预支付单
	var rec receive.PaySuccessReqStruct
	var err error
	rec.OrderId, _ = strconv.ParseInt(ctx.PostForm("orderId"), 10, 64)
	rec.PayType, _ = strconv.Atoi(ctx.PostForm("payType"))
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	//查一下订单的详情
	orderInfo := &order2.OmsOrder{}
	if err = global.Db.Model(&order2.OmsOrder{}).Where("id = ?", rec.OrderId).Find(&orderInfo).Error; err != nil {
		global.Logger.Error("支付订单时，根据订单id查询订单详情failed:%v", err)
	}
	//1.//生成本地唯一流水单号
	paymentSn := uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	payment := &order2.OmsPayment{}
	if err := global.Db.Model(&order2.OmsPayment{}).Debug().Where("order_id = ? and pay_status = ?", rec.OrderId, 0).Find(&payment).Error; err != nil {
		global.Logger.Errorf("查询重复订单流水号错误：%v", err)
		c.Response(ctx, "数据库繁忙，请稍后再试", 0, nil)
	}
	if payment.Id != 0 {
		global.Logger.Errorf("请勿重复支付，订单id=%d", rec.OrderId)
		c.Response(ctx, "请勿重复支付", 0, nil)
	}
	//然后在oms_payment表插入对应的项
	payment.OutTradeNo = paymentSn
	payment.UserId = memberId
	payment.PayStatus = 1 //正在支付
	payment.PayType = strconv.Itoa(rec.PayType)
	payment.OrderId = rec.OrderId
	payment.TotalAmount = float64(orderInfo.PayAmount)

	if err := payment.Insert(); err != nil {
		global.Logger.Errorf("支付流水号插入失败：%d", rec.OrderId)
		c.Response(ctx, "支付流水号插入失败", 0, nil)
		return
	}
	//2.todo:这里应该请求一下本地的 localhost:9090/alipay/pay接口，将订单sn号和订单金额作为参数传递过去
	//已经生成了预支付订单，这里就应该执行alipay/pay的逻辑了

	count, err := order.PaySuccess(&rec, memberId)
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
