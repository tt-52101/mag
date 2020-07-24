// @Title: router.go
// @Author: key7men@gmail.com
// @Description: 服务端API接口定义
// @Update: 2020/7/22 4:07 PM 
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
	Auth           	auth.Auther
	CasbinEnforcer 	*casbin.SyncedEnforcer
	DemoAPI        	*handler.Demo
	LoginAPI 	   	*handler.Login
	MenuAPI 		*handler.Menu
	RoleAPI 		*handler.Role
	UserAPI			*handler.User
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

// RegisterAPI register api group router
func (r *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	g.Use(middleware.UserAuthMiddleware(r.Auth,
		middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
	))

	g.Use(middleware.CasbinMiddleware(r.CasbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/pub"),
	))

	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.GET("captchaid", r.LoginAPI.GetCaptchaId)
				gLogin.GET("captcha", r.LoginAPI.GetCaptchaPic)
				gLogin.POST("", r.LoginAPI.Login)
				gLogin.POST("exit", r.LoginAPI.Logout)
			}

			gCurrent := pub.Group("current")
			{
				gCurrent.PUT("password", r.LoginAPI.UpdatePassword)
				gCurrent.GET("user", r.LoginAPI.GetUserInfo)
				gCurrent.GET("menutree", r.LoginAPI.QueryUserMenuTree)
			}
			pub.POST("/refresh-token", r.LoginAPI.RefreshToken)
		}

		gDemo := v1.Group("demos")
		{
			gDemo.GET("", r.DemoAPI.Query)
			gDemo.GET(":id", r.DemoAPI.Get)
			gDemo.POST("", r.DemoAPI.Create)
			gDemo.PUT(":id", r.DemoAPI.Update)
			gDemo.DELETE(":id", r.DemoAPI.Delete)
			gDemo.PATCH(":id/enable", r.DemoAPI.Enable)
			gDemo.PATCH(":id/disable", r.DemoAPI.Disable)
		}

		gMenu := v1.Group("menus")
		{
			gMenu.GET("", r.MenuAPI.Query)
			gMenu.GET(":id", r.MenuAPI.Get)
			gMenu.POST("", r.MenuAPI.Create)
			gMenu.PUT(":id", r.MenuAPI.Update)
			gMenu.DELETE(":id", r.MenuAPI.Delete)
			gMenu.PATCH(":id/enable", r.MenuAPI.Enable)
			gMenu.PATCH(":id/disable", r.MenuAPI.Disable)
		}
		v1.GET("/menus.tree", r.MenuAPI.QueryTree)

		gRole := v1.Group("roles")
		{
			gRole.GET("", r.RoleAPI.Query)
			gRole.GET(":id", r.RoleAPI.Get)
			gRole.POST("", r.RoleAPI.Create)
			gRole.PUT(":id", r.RoleAPI.Update)
			gRole.DELETE(":id", r.RoleAPI.Delete)
			gRole.PATCH(":id/enable", r.RoleAPI.Enable)
			gRole.PATCH(":id/disable", r.RoleAPI.Disable)
		}
		v1.GET("/roles.select", r.RoleAPI.QuerySelect)

		gUser := v1.Group("users")
		{
			gUser.GET("", r.UserAPI.Query)
			gUser.GET(":id", r.UserAPI.Get)
			gUser.POST("", r.UserAPI.Create)
			gUser.PUT(":id", r.UserAPI.Update)
			gUser.DELETE(":id", r.UserAPI.Delete)
			gUser.PATCH(":id/enable", r.UserAPI.Enable)
			gUser.PATCH(":id/disable", r.UserAPI.Disable)
		}
	}
}
