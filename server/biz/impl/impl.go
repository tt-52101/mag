package impl

import "github.com/google/wire"

// BizImplSet 注入
var BizImplSet = wire.NewSet(
	DemoSet,
	)