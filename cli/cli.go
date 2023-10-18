package cli

import (
	"raiton/cli/repl"

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
		Commands: []*cli.Command{
			{
				Name:   "repl",
				Usage:  "start the REPL",
				Action: repl.Run,
			},
			{
				Name:      "tokenize",
				Usage:     "tokenize the given file",
				ArgsUsage: "[file path]",
				Action:    tokenize,
			},
			{
				Name:      "parse",
				Usage:     "parse the given file and check for errors",
				ArgsUsage: "[file path]",
				Action:    parse,
			},
		},
	}

	return Cli{app: app}
}

func (c *Cli) Run(args []string) error {
	return c.app.Run(args)
}
