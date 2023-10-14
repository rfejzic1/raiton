package cli

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/parser"
	"github.com/urfave/cli/v2"
)

func parse(ctx *cli.Context) error {
	filePath := ctx.Args().First()

	if filePath == "" {
		return cli.Exit("expected a path to file for parsing", 1)
	}

	filePath = path.Clean(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	source, err := io.ReadAll(f)

	l := lexer.New(string(source))
	p := parser.New(&l)

	program, err := p.Parse()

	if err != nil {
		fmt.Fprintf(ctx.App.ErrWriter, "parse error: %s\n", err)
	}

	printer := ast.NewPrinter(program)

	fmt.Fprintln(ctx.App.Writer, printer.String())

	return nil
}
