package auth

import (
	"context"
	"errors"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

// Auther 认证接口
type Auther interface {
	// 生成令牌
	GenerateToken(ctx context.Context, userID string) (TokenInfo, error)

	// 销毁令牌
	DestroyToken(ctx context.Context, accessToken string) error

	// 解析用户ID
	ParseUserID(ctx context.Context, accessToken string) (string, error)

	// 释放资源
	Release() error
}
