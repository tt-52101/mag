package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/key7men/mag/server/biz"
	egin "github.com/key7men/mag/server/enhance/gin"
	"github.com/key7men/mag/server/schema"
)

// DemoSet 注入Demo
var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"))

// Demo 示例程序
type Demo struct {
	DemoBiz biz.IDemo
}

// Query 查询数据
func (a *Demo) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.DemoQueryParam
	if err := egin.ParseQuery(c, &params); err != nil {
		egin.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.DemoBiz.Query(ctx, params)
	if err != nil {
		egin.ResError(c, err)
		return
	}

	egin.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
func (a *Demo) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.DemoBiz.Get(ctx, c.Param("id"))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, item)
}

// Create 创建数据
func (a *Demo) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Demo
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	item.Creator = egin.GetUserID(c)
	result, err := a.DemoBiz.Create(ctx, item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, result)
}

// Update 更新数据
func (a *Demo) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Demo
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	err := a.DemoBiz.Update(ctx, c.Param("id"), item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Delete 删除数据
func (a *Demo) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.DemoBiz.Delete(ctx, c.Param("id"))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Enable 启用数据
func (a *Demo) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.DemoBiz.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Disable 禁用数据
func (a *Demo) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.DemoBiz.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}
