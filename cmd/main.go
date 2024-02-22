package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/serhii-marchuk/blog/internal/bootstrap"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
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
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"c"},
				Usage:   "Count of the migration ",
			},
		},
		Action: func(cCtx *cli.Context) error {
			l := web.NewAppLogger()
			cfg := configs.NewDbConfig(l)
			d := cCtx.String("direction")
			bootstrap.NewMigrator(d).RunDbMigration(cfg, l)

			return nil
		},
	},
	{
		Name:  "web:start",
		Usage: "start web server",
		Action: func(cCtx *cli.Context) error {
			fx.New(
				fx.Provide(web.NewAppLogger),
				fx.Provide(bootstrap.NewDb),
				fx.Provide(web.NewWebServer),
				fx.Provide(web.NewRenderer),
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

	if godotenv.Load("./configs/.env") != nil {
		os.Exit(0)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Something went wrong with app!")
	}
}
