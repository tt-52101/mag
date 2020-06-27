package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/key7men/mag/server/config"
	"github.com/key7men/mag/server/middleware"
	"github.com/key7men/mag/server/assist/gzip"
	"github.com/key7men/mag/server/router"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitGinEngine 初始化gin引擎
func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	prefixes := r.Prefixes()

	// Trace ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Copy body
	app.Use(middleware.CopyBodyMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Access logger
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Recover
	app.Use(middleware.RecoveryMiddleware())

	// CORS
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// GZIP
	if config.C.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExts(config.C.GZIP.ExcludedExts),
			gzip.WithExcludedPaths(config.C.GZIP.ExcludedPaths),
		))
	}

	// Router register
	r.Register(app)

	// Swagger
	if config.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Static Website
	if dir := config.C.Static; dir != "" {
		app.Use(middleware.StaticMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	return app
}
