package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/key7men/mag/pkg/auth"
	"github.com/key7men/mag/pkg/errs"
	"github.com/key7men/mag/pkg/logger"
	"github.com/key7men/mag/server/config"
	icontext "github.com/key7men/mag/server/enhance/context"

	egin "github.com/key7men/mag/server/enhance/gin"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	egin.SetUserID(c, userID)
	ctx := icontext.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), egin.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, config.C.Root.UserName)
					c.Next()
					return
				}
				egin.ResError(c, errs.ErrInvalidToken)
				return
			}
			egin.ResError(c, errs.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
