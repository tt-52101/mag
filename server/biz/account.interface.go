package biz

import (
	"context"

	"github.com/key7men/mag/server/schema"
)

// ILogin 登录业务逻辑接口
type ILogin interface {
	// 登录验证
	Verify(ctx context.Context, userName, password string) (*schema.User, error)
	// 生成令牌
	GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error)
	// 销毁令牌
	DestroyToken(ctx context.Context, tokenString string) error
	// 获取用户登录信息
	GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error)
	// 查询用户的权限菜单树
	QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error)
	// 更新用户登录密码
	UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error
}
