// @Title: account.handler.go
// @Author: key7men@gmail.com
// @Description: 账户登录登出接口处理
// @Update: 2020/7/23 4:22 PM 
package handler

import (
	"github.com/dchest/captcha"
	"github.com/key7men/mag/pkg/errs"
	"github.com/key7men/mag/server/biz"
	"github.com/key7men/mag/server/config"
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
	LoginBiz biz.ILogin
}

// GetCaptchaId 获取验证码ID
func (l *Login) GetCaptchaId(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := l.LoginBiz.GetCaptchaId(ctx, config.C.Captcha.Length)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, item)
}

// GetCaptchaPic 获取验证码图片
func (l *Login) GetCaptchaPic(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		egin.ResError(c,errs.New400Response("请提供验证码ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			egin.ResError(c,errs.New400Response("未找到验证码ID"))
			return
		}
	}

	cfg := config.C.Captcha
	err := l.LoginBiz.GetCaptchaPic(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		egin.ResError(c,err)
	}
}

// Login 用户登录
func (l *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}


	user, err := l.LoginBiz.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		egin.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	egin.SetUserID(c, userID)

	ctx = logger.NewUserIDContext(ctx, userID)
	tokenInfo, err := l.LoginBiz.GenerateToken(ctx, userID)
	if err != nil {
		egin.ResError(c, err)
		return
	}

	logger.StartSpan(ctx, logger.SetSpanTitle("用户登录"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	egin.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
func (l *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := egin.GetUserID(c)
	if userID != "" {
		err := l.LoginBiz.DestroyToken(ctx, egin.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ctx, logger.SetSpanTitle("用户登出"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	egin.ResOK(c)
}

// RefreshToken 刷新令牌
func (l *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := l.LoginBiz.GenerateToken(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
func (l *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := l.LoginBiz.GetLoginInfo(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
func (l *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := l.LoginBiz.QueryUserMenuTree(ctx, egin.GetUserID(c))
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResList(c, menus)
}

// UpdatePassword 更新个人密码
func (l *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := egin.ParseJSON(c, &item); err != nil {
		egin.ResError(c, err)
		return
	}

	err := l.LoginBiz.UpdatePassword(ctx, egin.GetUserID(c), item)
	if err != nil {
		egin.ResError(c, err)
		return
	}
	egin.ResOK(c)
}
