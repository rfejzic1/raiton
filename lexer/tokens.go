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
	"(":  LEFT_PAREN,
	")":  RIGHT_PAREN,
	"[":  LEFT_BRACKET,
	"]":  RIGHT_BRACKET,
	"{":  LEFT_BRACE,
	"}":  RIGHT_BRACE,
	"<":  LEFT_ANGLE,
	">":  RIGHT_ANGLE,
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

	LEFT_PAREN    = "left_paren"
	RIGHT_PAREN   = "right_paren"
	LEFT_BRACKET  = "left_bracket"
	RIGHT_BRACKET = "right_bracket"
	LEFT_BRACE    = "left_brace"
	RIGHT_BRACE   = "right_brace"
	LEFT_ANGLE    = "left_angle"
	RIGHT_ANGLE   = "right_angle"

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
