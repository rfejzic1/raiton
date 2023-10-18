package lexer

import (
	"testing"

	"raiton/token"
)

type test struct {
	name string
	t    *testing.T
}

type tokenExpect struct {
	Type    token.TokenType
	Literal string
}

func newTest(t *testing.T, name string) test {
	return test{t: t, name: name}
}

func (t *test) expect(source string, sequence []tokenExpect) {
	l := New(source)

	for i, et := range sequence {
		token := l.Next()

		if token.Type != et.Type {
			t.t.Fatalf("%s[%d] - wrong token type; expected `%s`, but got `%s` at line %d", t.name, i, et.Type, token.Type, token.Line)
		}

		if token.Literal != et.Literal {
			t.t.Fatalf("%s[%d] - wrong token literal; expected `%s`, but got `%s` at line %d", t.name, i, string(et.Literal), string(token.Literal), token.Line)
		}
	}
}

func TestEmptySource(t *testing.T) {
	test := newTest(t, "TestEmptySource")
	source := ``

	test.expect(source, []tokenExpect{})
}

func TestIdentifierLexing(t *testing.T) {
	test := newTest(t, "IdentifierLexing")
	source := `
	add
	println!
	_
	_what
	___
	_what_is_this
	_1
	_a1
	_1a
	_!
	`

	test.expect(source, []tokenExpect{
		{token.IDENTIFIER, `add`},
		{token.IDENTIFIER, `println!`},
		{token.IDENTIFIER, `_`},
		{token.IDENTIFIER, `_what`},
		{token.IDENTIFIER, `___`},
		{token.IDENTIFIER, `_what_is_this`},
		{token.IDENTIFIER, `_1`},
		{token.IDENTIFIER, `_a1`},
		{token.IDENTIFIER, `_1a`},
		{token.IDENTIFIER, `_!`},
		{token.EOF, ``},
	})
}

func TestNumberLexing(t *testing.T) {
	test := newTest(t, "NumberLexing")
	source := `
	1
	0.1
	-0.1
	-0.
	-2
	1.
	`

	test.expect(source, []tokenExpect{
		{token.NUMBER, `1`},
		{token.NUMBER, `0`},
		{token.DOT, `.`},
		{token.NUMBER, `1`},
		{token.MINUS, `-`},
		{token.NUMBER, `0`},
		{token.DOT, `.`},
		{token.NUMBER, `1`},
		{token.MINUS, `-`},
		{token.NUMBER, `0`},
		{token.DOT, `.`},
		{token.MINUS, `-`},
		{token.NUMBER, `2`},
		{token.NUMBER, `1`},
		{token.DOT, `.`},
		{token.EOF, ``},
	})
}

func TestStringLexing(t *testing.T) {
	test := newTest(t, "StringLexing")
	source := `"Hello, Raiton!"`

	test.expect(source, []tokenExpect{
		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `Hello, Raiton!`},
		{token.DOUBLE_QUOTE, `"`},
		{token.EOF, ``},
	})
}

func TestSkippingSpaces(t *testing.T) {
	test := newTest(t, "TestSkippingSpaces")
	source := `  println  123.1 "Raiton"  `

	test.expect(source, []tokenExpect{
		{token.IDENTIFIER, `println`},
		{token.NUMBER, `123`},
		{token.DOT, `.`},
		{token.NUMBER, `1`},
		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `Raiton`},
		{token.DOUBLE_QUOTE, `"`},
		{token.EOF, ``},
	})
}

func TestSkippingNewlines(t *testing.T) {
	test := newTest(t, "TestSkippingNewlines")
	source := `  println
	123.1 
	   "Raiton"  
	`

	test.expect(source, []tokenExpect{
		{token.IDENTIFIER, `println`},
		{token.NUMBER, `123`},
		{token.DOT, `.`},
		{token.NUMBER, `1`},
		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `Raiton`},
		{token.DOUBLE_QUOTE, `"`},
		{token.EOF, ``},
	})
}

func TestSkippingComments(t *testing.T) {
	test := newTest(t, "TestSkippingComments")
	source := `
	# comment 1
	ident # comment 2
	123  
	   # comment 3
	   "Raiton"  
	# comment 4
	3.14
	`

	test.expect(source, []tokenExpect{
		{token.IDENTIFIER, `ident`},
		{token.NUMBER, `123`},
		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `Raiton`},
		{token.DOUBLE_QUOTE, `"`},
		{token.NUMBER, `3`},
		{token.DOT, `.`},
		{token.NUMBER, `14`},
		{token.EOF, ``},
	})
}

func TestParenBracketBraceAngleLexing(t *testing.T) {
	test := newTest(t, "TestParenBracketBraceAngleLexing")
	source := `()[]{}`

	test.expect(source, []tokenExpect{
		{token.OPEN_PAREN, `(`},
		{token.CLOSED_PAREN, `)`},
		{token.OPEN_BRACKET, `[`},
		{token.CLOSED_BRACKET, `]`},
		{token.OPEN_BRACE, `{`},
		{token.CLOSED_BRACE, `}`},
		{token.EOF, ``},
	})
}

func TestQuoteLexing(t *testing.T) {
	test := newTest(t, "TestQuoteLexing")
	source := `
	# should parse quotes correctly
	"'single'"
	'"double"'
	"double escape \""
	'single escape \''
	`

	test.expect(source, []tokenExpect{
		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `'single'`},
		{token.DOUBLE_QUOTE, `"`},

		{token.SINGLE_QUOTE, `'`},
		{token.STRING, `"double"`},
		{token.SINGLE_QUOTE, `'`},

		{token.DOUBLE_QUOTE, `"`},
		{token.STRING, `double escape "`},
		{token.DOUBLE_QUOTE, `"`},

		{token.SINGLE_QUOTE, `'`},
		{token.STRING, `single escape '`},
		{token.SINGLE_QUOTE, `'`},
	})
}

func TestFunctionLexing(t *testing.T) {
	test := newTest(t, "TestFunctionLexing")
	source := `
	(map [1 2 3] \x: (square x))
	`

	test.expect(source, []tokenExpect{
		{token.OPEN_PAREN, `(`},
		{token.IDENTIFIER, `map`},
		{token.OPEN_BRACKET, `[`},
		{token.NUMBER, `1`},
		{token.NUMBER, `2`},
		{token.NUMBER, `3`},
		{token.CLOSED_BRACKET, `]`},
		{token.BACKSLASH, `\`},
		{token.IDENTIFIER, `x`},
		{token.COLON, `:`},
		{token.OPEN_PAREN, `(`},
		{token.IDENTIFIER, `square`},
		{token.IDENTIFIER, `x`},
		{token.CLOSED_PAREN, `)`},
		{token.CLOSED_PAREN, `)`},
	})
}
