package lexer

import "fmt"

type TokenType string

type Token struct {
	Literal string
	Line    int
	Column  int
	Type    TokenType
}

func (t *Token) Print() {
	fmt.Printf("(%3d, %3d) %12s '%s'\n", t.Line, t.Column, t.Type, t.Literal)
}

var SYMBOLS = map[string]TokenType{
	"(":  LEFT_PAREN,
	")":  RIGHT_PAREN,
	":":  COLON,
}

const (
	IDENTIFIER    = "identifier"
	STRING        = "string"
	NUMBER        = "number"
	LEFT_PAREN    = "left_paren"
	RIGHT_PAREN   = "right_paren"
	COLON         = "colon"

	EOF     = "eof"
	UNKNOWN = "unknown"
)
