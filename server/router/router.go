package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/key7men/mag/pkg/auth"
	"github.com/key7men/mag/server/handler"
	"github.com/key7men/mag/server/middleware"
)

var _ IRouter = (*Router)(nil)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	DemoAPI        *handler.Demo
}

// Register 注册路由
func (r *Router) Register(app *gin.Engine) error {
	r.RegisterAPI(app)
	return nil
}

// Prefixes 路由前缀列表
func (r *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

func (r *Router) RegisterAPI(app *gin.Engine) {
	api := app.Group("/api")

	// 加入用户认证中间件
	api.Use(middleware.UserAuthMiddleware(r.Auth,
		middleware.AllowMethodAndPathPrefixSkipper("/api/v1/pub/login"),
	))

	// 加入casbin中间件
	api.Use(middleware.CasbinMiddleware(r.CasbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/demo"),
	))

	api.Use(middleware.RateLimiterMiddleware())

	v1 := api.Group("/v1")
	{
		apiDemo := v1.Group("demo")
		{
			apiDemo.GET("", r.DemoAPI.Query)
			apiDemo.GET(":id", r.DemoAPI.Get)
		}
	}
}
