package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/rfejzic1/raiton/evaluator"
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

		node, err := par.Parse()

		if err != nil {
			fmt.Fprintf(out, "error: %s\n", err)
			continue
		}

		eval := evaluator.New(node)

		obj, err := eval.Evaluate()

		if err != nil {
			fmt.Fprintf(out, "error: %s\n", err)
			continue
		}

		fmt.Fprintf(out, "%s\n", obj.Inspect())
	}

	return nil
}
