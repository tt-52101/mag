package provider

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/key7men/mag/pkg/auth"
)

type Provider struct {
	Engine 		*gin.Engine
	Auth		auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
}

var ProviderSet = wire.NewSet(wire.Struct(new(Provider), "*"))