package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/key7men/mag/pkg/errs"
	"github.com/key7men/mag/server/config"
	egin "github.com/key7men/mag/server/enhance/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.Enforce(egin.GetUserID(c), p, m); err != nil {
			egin.ResError(c, errs.WithStack(err))
			return
		} else if !b {
			egin.ResError(c, errs.ErrNoPerm)
			return
		}
		c.Next()
	}
}
