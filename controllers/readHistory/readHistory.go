package readHistory

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/readHistory"
	"gomall/logic/readHistory"
	"gomall/utils/jwt"
)

type ReadHistoryController struct {
	controller.BaseControllers
}

func (c *ReadHistoryController) Create(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.CreateReadHistoryReqStruct)); err == nil {
		//global.Logger.Infof("创建浏览记录的传参:%v", rec)
		memberId, err := jwt.GetMemberIdFromCtx(ctx)
		if err != nil {
			c.Response(ctx, "用户身份校验错误", nil, err)
		}
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
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("获取浏览记录绑定参数失败:%v", err)
		c.Response(ctx, "请求参数错误", nil, nil)
		return
	}
	memberIdFromCtx, _ := jwt.GetMemberIdFromCtx(ctx)

	readHistoryList, err := readHistory.ListReadHistory((rec.PageNum-1)*rec.PageSize, rec.PageSize, memberIdFromCtx)
	if err != nil {
		c.Response(ctx, "分页获取浏览记录出错", nil, err)
		return
	}
	c.Response(ctx, "分页获取浏览记录成功", readHistoryList, nil)
}
