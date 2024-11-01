package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/coupon"
	cartLogic "gomall/logic/cart"
	coupon "gomall/logic/coupon"
	"gomall/utils/jwt"
)

type CouponController struct {
	controller.BaseControllers
}

func (c *CouponController) Add(ctx *gin.Context) {
	var rec receive.AddCouponRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("领取优惠券时，参数绑定失败:%v", err.Error())
		c.Response(ctx, "领取优惠券时，参数绑定失败:"+err.Error(), nil, err)
		return
	}
	if rec.CouponId == 0 {
		c.Response(ctx, "领取优惠券时，优惠券id非法", nil, errors.New("领取优惠券时，优惠券id非法"))
	}

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	if err := coupon.AddCoupon(rec.CouponId, memberId); err != nil {
		c.Response(ctx, "领取优惠券时，执行领取操作失败:"+err.Error(), nil, err)
		return
	}

	c.Response(ctx, "领取优惠券成功", nil, nil)
}

func (c *CouponController) List(ctx *gin.Context) {
	var rec receive.ListCouponRequestStruct
	//请求优惠券列表的参数是字符串参数，用shouldBind绑定
	if err := ctx.ShouldBind(&rec); err != nil {
		global.Logger.Errorf("获取优惠券列表时，参数绑定失败:%v", err.Error())
		c.Response(ctx, "获取优惠券列表时，参数绑定失败:"+err.Error(), nil, err)
		return
	}
	global.Logger.Infof("接收到的coupon的useStatus为：%d", rec.UseStatus)
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	couponList, err := coupon.GetCouponList(memberId, rec.UseStatus)
	if err != nil {
		c.Response(ctx, "获取优惠券列表失败:"+err.Error(), nil, err)
		return
	}

	c.Response(ctx, "获取优惠券列表成功", couponList, nil)
}

func (c *CouponController) ListCart(ctx *gin.Context) {
	var rec receive.ListCartRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("获取购物车的相关优惠券时，参数绑定失败:%v", err.Error())
		c.Response(ctx, "获取购物车的相关优惠券时，参数绑定失败:"+err.Error(), nil, err)
		return
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	cartPromotionItemList, err := cartLogic.CartListPromotion(nil, memberId)
	if err != nil {
		c.Response(ctx, "获取购物车的相关优惠券时，获取购物车商品列表失败", nil, err)
		return
	}
	couponHistoryList, err := coupon.ListCartCoupon(cartPromotionItemList, rec.Type, memberId)
	if err != nil {
		c.Response(ctx, "获取购物车的相关优惠券时,查询购物车商品对应的优惠券失败:"+err.Error(), nil, err)
		return
	}
	c.Response(ctx, "获取购物车的相关优惠券成功", couponHistoryList, nil)
}

// ListByProduct 获取当前商品相关的优惠券
func (c *CouponController) ListByProduct(ctx *gin.Context) {
	var rec receive.ListByProductRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("获取当前商品相关的优惠券时，参数绑定失败:%v", err.Error())
		c.Response(ctx, "获取当前商品相关的优惠券时，参数绑定失败:"+err.Error(), nil, err)
		return
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	couponHistoryList, err := coupon.ListByProductId(rec.ProductId, memberId)
	if err != nil {
		c.Response(ctx, "获取当前商品相关的优惠券失败:"+err.Error(), nil, err)
		return
	}
	c.Response(ctx, "获取当前商品相关的优惠券成功", couponHistoryList, nil)
}

// 获取会员优惠券历史列表
func (c *CouponController) ListHistory(ctx *gin.Context) {
	var rec receive.ListHistoryRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("获取会员优惠券历史列表时，参数绑定失败:%v", err.Error())
		c.Response(ctx, "获取会员优惠券历史列表时，参数绑定失败:"+err.Error(), nil, err)
		return
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	couponHistoryList, err := coupon.ListHistory(memberId, rec.UseStatus)
	if err != nil {
		c.Response(ctx, "获取会员优惠券历史列表失败:"+err.Error(), nil, err)
		return
	}
	c.Response(ctx, "获取会员优惠券历史列表成功", couponHistoryList, nil)
}
