package main

import (
	"fmt"
	"github.com/serhii-marchuk/blog/internal/bootstrap"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/constructors"
	"github.com/serhii-marchuk/blog/internal/bootstrap/web"
	webHandl "github.com/serhii-marchuk/blog/internal/ports/web"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"os"
)

var commands = []*cli.Command{
	{
		Name:  "migrate",
		Usage: "options for migrates",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "direction",
				Aliases: []string{"d"},
				Usage:   "Point the direction of migration",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cfg := configs.NewConfigs()
			d := cCtx.String("direction")
			bootstrap.NewMigrator(d).RunDbMigration(&cfg.Database)

			return nil
		},
	},
	{
		Name:  "web:start",
		Usage: "start web server",
		Action: func(cCtx *cli.Context) error {
			fx.New(
				fx.Provide(func() *configs.Configs {
					return configs.NewConfigs()
				}),
				fx.Provide(func() *configs.WebCfg {
					return configs.NewWebCfg()
				}),
				fx.Provide(constructors.NewLogger),
				fx.Provide(constructors.NewRedisClient),
				fx.Provide(constructors.NewDb),
				fx.Provide(constructors.NewRenderer),
				fx.Provide(web.NewWebServer),
				fx.Provide(webHandl.NewWebHandler),
				fx.Invoke(
					web.Setup,
					web.Start,
				),
			).Run()
			return nil
		},
	},
}

func main() {
	app := &cli.App{Commands: commands}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Something went wrong with app!")
	}
}
