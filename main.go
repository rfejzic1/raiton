package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/parser"
)

const VERSION = "v0.0.1"

func main() {
	fmt.Printf("Raiton %s\n", VERSION)

	for {
		fmt.Print("> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			break
		}

		lex := lexer.New(input)
		par := parser.New(&lex)

		_, err := par.Parse()

		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Println("ok")
		}
	}
}
