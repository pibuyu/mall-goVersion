package users

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/productCollection"
	productCollectionLogic "gomall/logic/productCollection"
	productCollectionModels "gomall/models/productCollection"
	"gomall/models/users"
	"gomall/utils/jwt"
	"gomall/utils/response"
)

type ProductCollectionController struct {
	controller.BaseControllers
}

func (c *ProductCollectionController) Add(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.AddReqStruct)); err == nil {
		//先预处理一下，把里面的用户信息赋初值
		curUser := &users.User{}
		memberId, _ := jwt.GetMemberIdFromCtx(ctx)
		if err := curUser.GetMemberById(memberId); err != nil {
			c.Response(ctx, "获取用户信息失败", 0, err)
		}
		rec.MemberId = memberId
		rec.MemberNickname = curUser.Nickname
		rec.MemberIcon = curUser.Icon

		count, err := productCollectionLogic.Add(rec)
		if err != nil {
			c.Response(ctx, "添加商品到收藏夹失败", 0, err)
			return
		}
		c.Response(ctx, "添加商品到收藏夹成功", count, nil)
	}
}

func (c *ProductCollectionController) Clear(ctx *gin.Context) {

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)

	err := productCollectionLogic.Clear(memberId)
	if err != nil {
		c.Response(ctx, "清空当前用户商品收藏列表失败", nil, err)
		return
	}
	c.Response(ctx, "清空当前用户商品收藏列表成功", nil, nil)
}
func (c *ProductCollectionController) Delete(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteRespStruct)); err == nil {
		memberId, _ := jwt.GetMemberIdFromCtx(ctx)
		count, err := productCollectionLogic.Delete(rec.ProductId, memberId)
		if err != nil {
			c.Response(ctx, "根据id删除商品收藏失败", 0, err)
			return
		}
		c.Response(ctx, "根据id删除商品收藏成功", count, nil)
	}
}

func (c *ProductCollectionController) Detail(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.DetailRespStruct)); err == nil {
		memberId, _ := jwt.GetMemberIdFromCtx(ctx)
		result, err := productCollectionLogic.Detail(rec.ProductId, memberId)
		if err != nil {
			c.Response(ctx, "根据id获取商品收藏详情失败", nil, err)
			return
		}
		c.Response(ctx, "根据id获取商品收藏详情成功", result, nil)
	}
}
func (c *ProductCollectionController) List(ctx *gin.Context) {
	var rec receive.ListReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		global.Logger.Errorf("请求收藏商品list时，参数绑定失败:%v", err)
		c.Response(ctx, "获取商品收藏列表时，参数绑定失败", nil, err)
		return
	}

	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	result, err := productCollectionLogic.List(rec.PageNum, rec.PageSize, memberId)
	if err != nil {
		c.Response(ctx, "获取商品收藏列表失败", nil, err)
		return
	}
	pageResult := response.ResetPage(result, int64(len(result)), rec.PageNum, rec.PageSize)
	if pageResult.List == nil {
		pageResult.List = []productCollectionModels.MemberProductCollection{}
	}
	c.Response(ctx, "获取商品收藏列表成功", pageResult, nil)
}
