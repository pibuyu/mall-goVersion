package es

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/es"
	"gomall/logic/es"
)

type EsProductController struct {
	controller.BaseControllers
}

// 导入所有数据库中商品到ES
func (c *EsProductController) ImportAll(ctx *gin.Context) {
	count, err := es.ImportAll()
	if err != nil {
		c.Response(ctx, "导入所有商品失败", 0, err)
		return
	}
	c.Response(ctx, "导入所有商品成功", count, nil)
}

func (c *EsProductController) Delete(ctx *gin.Context) {
	var rec receive.DeleteReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		c.Response(ctx, "删除es商品入参绑定失败", nil, err)
		global.Logger.Errorf("删除es商品入参绑定失败:%v", err)
		return
	}
	if err := es.DeleteById(rec.Id); err != nil {
		c.Response(ctx, "删除es商品失败", nil, err)
		global.Logger.Errorf("删除es商品失败:%v", err)
		return
	}
	c.Response(ctx, "删除es商品成功", nil, nil)
}

func (c *EsProductController) Create(ctx *gin.Context) {
	var rec receive.CreateReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		c.Response(ctx, "向es添加商品入参绑定失败", nil, err)
		global.Logger.Errorf("向es添加商品入参绑定失败:%v", err)
		return
	}
	if err := es.Create(rec.Id); err != nil {
		c.Response(ctx, "向es添加商品失败", nil, err)
		global.Logger.Errorf("向es添加商品失败:%v", err)
		return
	}
	c.Response(ctx, "向es添加商品成功", nil, nil)
}

func (c *EsProductController) SimpleSearch(ctx *gin.Context) {
	var rec receive.SimpleSearchReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		c.Response(ctx, "入参绑定失败", nil, err)
		global.Logger.Errorf("入参绑定失败:%v", err)
		return
	}
	result, err := es.SimpleSearch(&rec)
	if err != nil {
		c.Response(ctx, "在es中根据关键字简单搜索商品失败", nil, err)
		global.Logger.Errorf("在es中根据关键字简单搜索商品失败:%v", err)
		return
	}
	c.Response(ctx, "在es中根据关键字简单搜索商品成功", result, nil)

}
