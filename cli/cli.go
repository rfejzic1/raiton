package cli

import "github.com/urfave/cli/v2"

type Cli struct {
	app *cli.App
}

func New() Cli {
	app := &cli.App{
		Name:  "raiton",
		Usage: "the Raiton language toolchain",
		Action: func(ctx *cli.Context) error {
			return cli.ShowAppHelp(ctx)
		},
	}

	return Cli{app: app}
}

func (c *Cli) Run(args []string) error {
	return c.app.Run(args)
}
