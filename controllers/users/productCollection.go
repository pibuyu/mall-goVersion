package users

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	receive "gomall/interaction/receive/productCollection"
	productCollectionLogic "gomall/logic/productCollection"
	"gomall/models/users"
	"gomall/utils/jwt"
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
