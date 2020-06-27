package handler

import (
	"github.com/key7men/mag/server/biz"
	egin "github.com/key7men/mag/server/enhance/gin"
	"github.com/key7men/mag/server/schema"
	"github.com/key7men/mag/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

// Login 登录管理
type Login struct {
	LoginBll biz.ILogin
}

// Login 用户登录
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}


	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		egin.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	egin.SetUserID(c, userID)

	ctx = logger.NewUserIDContext(ctx, userID)
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, userID)
	if err != nil {
		egin.ResError(c, err)
		return
	}

	logger.StartSpan(ctx, logger.SetSpanTitle("用户登录"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	egin.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := egin.GetUserID(c)
	if userID != "" {
		err := a.LoginBll.DestroyToken(ctx, egin.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ctx, logger.SetSpanTitle("用户登出"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	egin.ResOK(c)
}

// RefreshToken 刷新令牌
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginBll.GetLoginInfo(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginBll.QueryUserMenuTree(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResList(c, menus)
}

// UpdatePassword 更新个人密码
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ctx, egin.GetUserID(c), item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}
