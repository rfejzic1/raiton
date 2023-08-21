package cli

import (
	"io"
	"os"
	"path"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/token"
	"github.com/urfave/cli/v2"
)

func tokenize(ctx *cli.Context) error {
	filePath := ctx.Args().First()

	if filePath == "" {
		return cli.Exit("expected a path to file for tokenization", 1)
	}

	filePath = path.Clean(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	source, err := io.ReadAll(f)

	l := lexer.New(string(source))

	for t := l.Next(); t.Type != token.EOF; t = l.Next() {
		t.Print(ctx.App.Writer)
	}

	return nil
}
