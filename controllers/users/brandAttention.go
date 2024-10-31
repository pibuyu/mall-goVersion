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
