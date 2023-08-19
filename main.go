package main

import (
	"log"
	"os"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/token"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: lang <path-to-file>")
	}

	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	l := lexer.New(string(data))

	for t := l.Next(); t.Matches(token.EOF); t = l.Next() {
		t.Print(os.Stdout)
	}
}
