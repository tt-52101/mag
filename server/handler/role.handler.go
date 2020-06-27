package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/key7men/mag/server/biz"
	egin "github.com/key7men/mag/server/enhance/gin"
	"github.com/key7men/mag/server/schema"
)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"))

// Role 角色管理
type Role struct {
	RoleBll biz.IRole
}

// Query 查询数据
func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := egin.ParseQuery(c, &params); err != nil {
		egin.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.RoleBll.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect 查询选择数据
func (a *Role) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := egin.ParseQuery(c, &params); err != nil {
		egin.ResError(c, err)
		return
	}

	result, err := a.RoleBll.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResList(c, result.Data)
}

// Get 查询指定数据
func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleBll.Get(ctx, c.Param("id"))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, item)
}

// Create 创建数据
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	item.Creator = egin.GetUserID(c)
	result, err := a.RoleBll.Create(ctx, item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, result)
}

// Update 更新数据
func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	err := a.RoleBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Delete 删除数据
func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.Delete(ctx, c.Param("id"))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Enable 启用数据
func (a *Role) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}

// Disable 禁用数据
func (a *Role) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}
