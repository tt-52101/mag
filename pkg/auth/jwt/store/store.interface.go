package store

import (
	"time"
)

// Storer 令牌存储接口
type Storer interface {
	// 存储令牌数据，并指定到期时间
	Set(tokenString string, expiration time.Duration) error
	// 检查令牌是否存在
	Check(tokenString string) (bool, error)
	// 关闭存储
	Close() error
}
