package users

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	receive "gomall/interaction/receive/brandAttention"
	"gomall/logic/brandAttention"
	"gomall/models/users"
	"gomall/utils/jwt"
	"time"
)

type BrandAttentionController struct {
	controller.BaseControllers
}

// 添加收藏品牌
func (c *BrandAttentionController) Add(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.AddBrandAttentionReqStruct)); err == nil {
		//先预处理一下，把里面的用户信息赋初值
		curUser := &users.User{}
		memberId, _ := jwt.GetMemberIdFromCtx(ctx)
		if err := curUser.GetMemberById(memberId); err != nil {
			c.Response(ctx, "添加收藏品牌时，获取用户信息失败", 0, err)
			return
		}
		rec.MemberId = memberId
		rec.MemberNickname = curUser.Nickname
		rec.MemberIcon = curUser.Icon
		rec.CreateTime = time.Now()
		count, err := brandAttention.Add(rec)
		if err != nil {
			c.Response(ctx, "添加收藏品牌失败", 0, err)
			return
		}
		c.Response(ctx, "添加收藏品牌成功", count, nil)
	}
}

// 清空用户关注品牌列表
func (c *BrandAttentionController) Clear(ctx *gin.Context) {
	memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
	if err != nil {
		c.Response(ctx, "清空用户关注品牌列表时，用户身份校验错误", nil, err)
	}
	if err := brandAttention.Clear(memberIdFromCtx); err != nil {
		c.Response(ctx, "清空用户关注品牌列表失败", nil, err)
	}
	c.Response(ctx, "清空用户关注品牌列表成功", nil, nil)
}

func (c *BrandAttentionController) Delete(ctx *gin.Context) {
	var rec receive.DeleteBrandAttentionReqStruct
	if err := ctx.ShouldBindJSON(&rec); err == nil {
		memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "删除用户关注品牌列表时，用户身份校验错误", 0, err)
			return
		}
		if err = brandAttention.Delete(rec.BrandId, memberIdFromCtx); err != nil {
			c.Response(ctx, "删除用户关注品牌失败", 0, err)
			return
		}
		c.Response(ctx, "删除用户关注品牌成功", 1, nil)
	}
}

func (c *BrandAttentionController) Detail(ctx *gin.Context) {
	var rec receive.DetailReqStruct
	if err := ctx.ShouldBindJSON(&rec); err == nil {
		memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "获取用户关注品牌的详情时，用户身份校验错误", 0, err)
		}
		result, err := brandAttention.Detail(rec.BrandId, memberIdFromCtx)
		if err != nil {
			c.Response(ctx, "获取用户关注品牌的详情失败", nil, err)
		}

		c.Response(ctx, "获取用户关注品牌的详情成功", result, nil)
	}
}

func (c *BrandAttentionController) List(ctx *gin.Context) {
	var rec receive.ListReqStruct
	if err := ctx.ShouldBindJSON(&rec); err == nil {
		memberIdFromCtx, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "分页获取用户关注的品牌时，用户身份校验错误", 0, err)
		}
		result, err := brandAttention.List(rec.PageNum, rec.PageSize, memberIdFromCtx)
		if err != nil {
			c.Response(ctx, "分页获取用户关注的品牌失败", nil, err)
		}

		c.Response(ctx, "分页获取用户关注的品牌成功", result, nil)
	}
}
