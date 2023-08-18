package main

import (
	"log"
	"os"

	"github.com/rfejzic1/raiton/lexer"
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

	for token := l.Next(); token.Type != lexer.EOF; token = l.Next() {
		token.Print()
	}
}
