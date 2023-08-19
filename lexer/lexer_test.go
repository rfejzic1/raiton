package lexer

import "testing"

type test struct {
	name string
	t    *testing.T
}

type token struct {
	Type    TokenType
	Literal string
}

func newTest(t *testing.T, name string) test {
	return test{t: t, name: name}
}

func (t *test) expect(source string, sequence []token) {
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

func TestIdentifierLexing(t *testing.T) {
	test := newTest(t, "IdentifierLexing")
	source := `println`

	test.expect(source, []token{
		{IDENTIFIER, `println`},
		{EOF, ``},
	})
}

func TestNumberLexing(t *testing.T) {
	test := newTest(t, "NumberLexing")
	source := `123`

	test.expect(source, []token{
		{NUMBER, `123`},
		{EOF, ``},
	})
}

func TestNumberLexingWithDecimal(t *testing.T) {
	test := newTest(t, "NumberLexingWithDecimal")
	source := `123.4`

	test.expect(source, []token{
		{NUMBER, `123.4`},
		{EOF, ``},
	})
}

func TestNumberLexingWithTrailingDecimal(t *testing.T) {
	test := newTest(t, "NumberLexingWithTrailingDecimal")
	source := `123.`

	test.expect(source, []token{
		{NUMBER, "123."},
		{EOF, ``},
	})
}

func TestStringLexing(t *testing.T) {
	test := newTest(t, "StringLexing")
	source := `"Hello, Raiton!"`

	test.expect(source, []token{
		{DOUBLE_QUOTE, `"`},
		{STRING, `Hello, Raiton!`},
		{DOUBLE_QUOTE, `"`},
		{EOF, ``},
	})
}

func TestSkippingSpaces(t *testing.T) {
	test := newTest(t, "TestSkippingSpaces")
	source := `  println  123.1 "Raiton"  `

	test.expect(source, []token{
		{IDENTIFIER, `println`},
		{NUMBER, `123.1`},
		{DOUBLE_QUOTE, `"`},
		{STRING, `Raiton`},
		{DOUBLE_QUOTE, `"`},
		{EOF, ``},
	})
}

func TestSkippingNewlines(t *testing.T) {
	test := newTest(t, "TestSkippingNewlines")
	source := `  println
	123.1 
	   "Raiton"  
	`

	test.expect(source, []token{
		{IDENTIFIER, `println`},
		{NUMBER, `123.1`},
		{DOUBLE_QUOTE, `"`},
		{STRING, `Raiton`},
		{DOUBLE_QUOTE, `"`},
		{EOF, ``},
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

	test.expect(source, []token{
		{IDENTIFIER, `ident`},
		{NUMBER, `123`},
		{DOUBLE_QUOTE, `"`},
		{STRING, `Raiton`},
		{DOUBLE_QUOTE, `"`},
		{NUMBER, `3.14`},
		{EOF, ``},
	})
}

func TestParenBracketBraceAngleLexing(t *testing.T) {
	test := newTest(t, "TestParenBracketBraceAngleLexing")
	source := `()[]{}<>`

	test.expect(source, []token{
		{LEFT_PAREN, `(`},
		{RIGHT_PAREN, `)`},
		{LEFT_BRACKET, `[`},
		{RIGHT_BRACKET, `]`},
		{LEFT_BRACE, `{`},
		{RIGHT_BRACE, `}`},
		{LEFT_ANGLE, `<`},
		{RIGHT_ANGLE, `>`},
		{EOF, ``},
	})
}
