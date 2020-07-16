package main

import (
	"context"
	"os"

	"github.com/key7men/mag/pkg/logger"
	"github.com/key7men/mag/server"
	"github.com/urfave/cli/v2"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.0"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTraceIDContext(context.Background(), "main") // main这里表明是main.go文件

	app := cli.NewApp()
	app.Name = "mag"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on GIN + GORM + REDIS + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "rbac",
				Aliases:  []string{"r"},
				Usage:    "casbin的访问控制模型(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Aliases: []string{"m"},
				Usage: "初始化菜单数据配置文件(.yaml)",
			},
			&cli.StringFlag{
				Name:  "www",
				Aliases:  []string{"w"},
				Usage: "静态站点目录",
			},
		},
		Action: func(c *cli.Context) error {
			return server.Run(ctx,
				server.SetConfigFile(c.String("conf")),
				server.SetModelFile(c.String("model")),
				server.SetStaticDir(c.String("www")),
				server.SetMenuFile(c.String("menu")),
				server.SetVersion(VERSION))
		},
	}
}
