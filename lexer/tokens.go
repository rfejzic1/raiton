package lexer

import (
	"fmt"
	"io"
)

type TokenType string

type Token struct {
	Literal string
	Line    int
	Column  int
	Type    TokenType
}

func (t *Token) Print(w io.Writer) {
	format := "(%3d, %3d) %12s %s\n"
	if t.Type == STRING {
		format = "(%3d, %3d) %12s `%s`\n"
	}
	fmt.Fprintf(w, format, t.Line, t.Column, t.Type, t.Literal)
}

var KEYWORDS = map[string]TokenType{
	"type": TYPE,
}

var SYMBOLS = map[string]TokenType{
	"(":  OPEN_PAREN,
	")":  CLOSED_PAREN,
	"[":  OPEN_BRACKET,
	"]":  CLOSED_BRACKET,
	"{":  OPEN_BRACE,
	"}":  CLOSED_BRACE,
	"<":  OPEN_ANGLE,
	">":  CLOSED_ANGLE,
	"'":  SINGLE_QUOTE,
	"\"": DOUBLE_QUOTE,
	".":  DOT,
	",":  COMMA,
	":":  COLON,
	"|":  PIPE,
	"\\": BACKSLASH,
	"->": RIGHT_ARROW,
}

const (
	IDENTIFIER = "identifier"
	STRING     = "string"
	NUMBER     = "number"

	TYPE = "type"

	OPEN_PAREN    = "left_paren"
	CLOSED_PAREN   = "right_paren"
	OPEN_BRACKET  = "left_bracket"
	CLOSED_BRACKET = "right_bracket"
	OPEN_BRACE    = "left_brace"
	CLOSED_BRACE   = "right_brace"
	OPEN_ANGLE    = "left_angle"
	CLOSED_ANGLE   = "right_angle"

	SINGLE_QUOTE = "single_quote"
	DOUBLE_QUOTE = "double_quote"
	DOT          = "dot"
	COMMA        = "comma"
	COLON        = "colon"
	PIPE         = "pipe"
	BACKSLASH    = "backslash"
	RIGHT_ARROW  = "right_arrow"

	EOF     = "eof"
	ILLEGAL = "illegal"
)
