package readHistory

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/readHistory"
	"gomall/logic/readHistory"
	readHistoryModels "gomall/models/history"
	"gomall/utils/jwt"
	"gomall/utils/response"
)

type ReadHistoryController struct {
	controller.BaseControllers
}

func (c *ReadHistoryController) Create(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.CreateReadHistoryReqStruct)); err == nil {
		memberId, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验错误", nil, err)
		}
		global.Logger.Infof("user %d try to create read history of product %d", memberId, rec.ProductId)
		if err := readHistory.CreateReadHistory(rec, memberId); err != nil {
			c.Response(ctx, "插入浏览记录出错", 0, err)
			return
		}
		c.Response(ctx, "插入浏览记录完成", 1, nil)
	}
}

func (c *ReadHistoryController) Clear(ctx *gin.Context) {
	memberId, err := jwt.GetMemberIdFromCtx(ctx)
	if err != nil {
		c.Response(ctx, "用户身份校验错误", nil, err)
	}
	if err := readHistory.ClearReadHistory(memberId); err != nil {
		c.Response(ctx, "清空浏览记录出错", nil, err)
	}
	c.Response(ctx, "清空浏览记录完成", nil, nil)
}

func (c *ReadHistoryController) Delete(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.DeleteReadHistoryReqStruct)); err == nil {
		memberId, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验错误", nil, err)
		}
		if err := readHistory.DeleteReadHistoryByIds(rec, memberId); err != nil {
			c.Response(ctx, "删除浏览记录出错", 0, err)
			return
		}
		c.Response(ctx, "删除浏览记录完成", 1, nil)
	}
}

func (c *ReadHistoryController) List(ctx *gin.Context) {
	var rec receive.ListReadHistoryReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		global.Logger.Errorf("获取浏览记录绑定参数失败:%v", err)
		c.Response(ctx, "请求参数错误", nil, nil)
		return
	}
	memberIdFromCtx, _ := jwt.GetMemberIdFromCtx(ctx)

	readHistoryList, err := readHistory.ListReadHistory(rec.PageNum, rec.PageSize, memberIdFromCtx)
	if err != nil {
		c.Response(ctx, "分页获取浏览记录出错", nil, err)
		return
	}
	pageResult := response.ResetPage(readHistoryList, int64(len(readHistoryList)), int(rec.PageNum), int(rec.PageSize))
	//确保返回的是空数组而不是nil
	if pageResult.List == nil {
		pageResult.List = []readHistoryModels.MemberReadHistoryMongoStruct{}
	}
	c.Response(ctx, "分页获取浏览记录成功", pageResult, nil)
}
