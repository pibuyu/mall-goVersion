package cart

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/cart"
	"gomall/logic/cart"
	"gomall/utils/jwt"
)

type CartController struct {
	controller.BaseControllers
}

// todo:原项目中应该通过返回值是0 or 1来判断操作是否成功，这里尊重原项目的写法
func (c *CartController) AddCartItem(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.AddCartItemRequestStruct)); err == nil {
		//在这里取出用户信息memberId
		token := ctx.Request.Header.Get("Authorization")
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.Response(ctx, "用户身份校验失败", 0, err)
			return
		}
		err = cart.AddCartItem(rec, claims)
		if err != nil {
			c.Response(ctx, "购物车添加商品失败", 0, err)
			return
		}
		c.Response(ctx, "购物车添加商品成功", 1, nil)
	}
}

func (c *CartController) Clear(ctx *gin.Context) {
	if _, err := controller.ShouldBind(ctx, new(receive.ClearCartRequestStruct)); err == nil {
		memberId, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验失败", 0, err)
			return
		}
		if err = cart.Clear(memberId); err != nil {
			c.Response(ctx, "清空购物车失败", 0, err)
			return
		}
		c.Response(ctx, "清空购物车成功", 1, err)
	}
}

func (c *CartController) DeleteByIds(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteCartItemsByIdsRequestStruct)); err == nil {
		memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验失败", 0, err)
			return
		}
		if err := cart.DeleteCartItemsByIds(memberIdFromCtx, rec.Ids); err != nil {
			c.Response(ctx, "删除购物车商品失败", 0, err)
			return
		}
		c.Response(ctx, "删除购物车商品成功", 1, err)
	}
}

// GetProductById 获取购物车中指定商品的规格,用于重选规格
// todo:传过来的参数是int类型时，controller.shouldBind方法绑定不上参数，是怎么回事？？？
func (c *CartController) GetProductById(ctx *gin.Context) {
	var rec receive.GetProductByIdRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("GetHotProductList请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}

	cartProduct, err := cart.GetProductById(&rec)
	if err != nil {
		c.Response(ctx, "获取购物车商品信息失败", 0, err)
		return
	}
	c.Response(ctx, "获取购物车商品信息成功", cartProduct, err)
}

func (c *CartController) List(ctx *gin.Context) {
	memberId, err := jwt.GetMemberIdFromCtx(ctx)
	if err != nil {
		c.Response(ctx, "用户身份校验失败", nil, err)
		return
	}
	results, err := cart.List(memberId)
	if err != nil {
		c.Response(ctx, "获取购物车全部信息失败", nil, err)
		return
	}
	c.Response(ctx, "获取购物车全部信息成功", results, err)
}

// CartListPromotion 获取当前会员的购物车列表，包括促销信息
func (c *CartController) CartListPromotion(ctx *gin.Context) {
	//参数没有绑定上
	var rec receive.CartListPromotionRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		c.Response(ctx, "获取当前会员的购物车列表时，参数绑定失败", nil, err)
		global.Logger.Errorf("获取当前会员的购物车列表时，参数绑定失败:%v", err)
		return
	}
	memberId, err := jwt.GetMemberIdFromCtx(ctx)
	if err != nil {
		c.Response(ctx, "获取当前会员的购物车列表时，用户身份校验失败", nil, err)
		return
	}
	cartPromotionItemList, err := cart.CartListPromotion(rec.CartIds, memberId)
	if err != nil {
		c.Response(ctx, "获取当前会员的购物车列表失败", nil, err)
		return
	}
	c.Response(ctx, "获取当前会员的购物车列表成功", cartPromotionItemList, nil)
}

// UpdateAttr 修改购物车中商品的规格
func (c *CartController) UpdateAttr(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UpdateAttrRequestStruct)); err == nil {
		if err := cart.UpdateAttr(rec); err != nil {
			c.Response(ctx, "修改购物车中商品的规格失败", 0, err)
			return
		}
		c.Response(ctx, "修改购物车中商品的规格成功", 1, err)
	}
}

// UpdateQuantity 修改购物车中指定商品商品的数量
func (c *CartController) UpdateQuantity(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UpdateQuantityRequestStruct)); err == nil {
		memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验失败", nil, err)
			return
		}
		if err := cart.UpdateQuantity(rec, memberIdFromCtx); err != nil {
			c.Response(ctx, "修改购物车中指定商品商品的数量失败", 0, err)
			return
		}
		c.Response(ctx, "修改购物车中指定商品商品的数量成功", 1, err)
	}
}
