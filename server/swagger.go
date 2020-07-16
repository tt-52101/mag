/*
Package app 生成swagger文档

文档规则请参考：https://github.com/swaggo/swag#declarative-comments-format

使用方式：

	go get -u github.com/swaggo/swag/cmd/swag
	swag init --generalInfo ./server/swagger.go --output ./server/swagger */
package server

// @title mag
// @version 1.0.0
// @description RBAC scaffolding based on GIN + GORM/MONGO + CASBIN + WIRE.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /

