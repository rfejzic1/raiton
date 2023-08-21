package cli

import (
	"github.com/urfave/cli/v2"
)

type Cli struct {
	app *cli.App
}

func New() Cli {
	app := &cli.App{
		Name:    "raiton",
		Usage:   "the Raiton language toolchain",
		Version: VERSION,
		Authors: []*cli.Author{
			{
				Name: "rfejzic1",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "repl",
				Usage:  "start the REPL",
				Action: repl,
			},
			{
				Name:  "tokenize",
				Usage: "tokenize the given file",
				ArgsUsage: "[file path]",
				Action:    tokenize,
			},
		},
	}

	return Cli{app: app}
}

func (c *Cli) Run(args []string) error {
	return c.app.Run(args)
}
