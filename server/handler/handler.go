package handler

import "github.com/google/wire"

// APISet 注入api
var HandlerSet = wire.NewSet(
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
