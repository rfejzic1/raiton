package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/parser"
	"github.com/urfave/cli/v2"
)

func repl(ctx *cli.Context) error {
	in := ctx.App.Reader
	out := ctx.App.Writer

	fmt.Fprintf(out, "Raiton %s\n", VERSION)

	for {
		fmt.Fprint(out, "> ")

		scanner := bufio.NewScanner(in)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			return err
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			break
		}

		lex := lexer.New(input)
		par := parser.New(&lex)

		_, err := par.Parse()

		if err != nil {
			fmt.Fprintf(out, "error: %s\n", err)
		} else {
			fmt.Fprintln(out, "ok")
		}
	}

	return nil
}
