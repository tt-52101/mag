// +build wireinject
// 必须在当前文件的第一行加上上述声明

package provider

import (
	"github.com/key7men/mag/server/handler"
	"github.com/key7men/mag/server/biz/impl"
	"github.com/key7men/mag/server/module/rbac"
	"github.com/key7men/mag/server/router"
	"github.com/google/wire"

	gormModel "github.com/key7men/mag/server/model/gorm/dao"
)

func BuildInjector() (*Provider, func(), error) {
	// 依赖顺序有wire去帮你构建，此处只需将你的provider加进来就行
	wire.Build(
		InitGormDB,
		gormModel.ModelSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		impl.BizImplSet,
		handler.HandlerSet,
		router.RouterSet,
		rbac.CasbinAdapterSet,
		ProviderSet,
	)
	return new(Provider), nil, nil // 本质上返回值没有任何含义
}