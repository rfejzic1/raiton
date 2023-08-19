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

func TestEmptySource(t *testing.T) {
	test := newTest(t, "TestEmptySource")
	source := ``

	test.expect(source, []token{})
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
		{OPEN_PAREN, `(`},
		{CLOSED_PAREN, `)`},
		{OPEN_BRACKET, `[`},
		{CLOSED_BRACKET, `]`},
		{OPEN_BRACE, `{`},
		{CLOSED_BRACE, `}`},
		{OPEN_ANGLE, `<`},
		{CLOSED_ANGLE, `>`},
		{EOF, ``},
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

	test.expect(source, []token{
		{DOUBLE_QUOTE, `"`},
		{STRING, `'single'`},
		{DOUBLE_QUOTE, `"`},

		{SINGLE_QUOTE, `'`},
		{STRING, `"double"`},
		{SINGLE_QUOTE, `'`},

		{DOUBLE_QUOTE, `"`},
		{STRING, `double escape "`},
		{DOUBLE_QUOTE, `"`},

		{SINGLE_QUOTE, `'`},
		{STRING, `single escape '`},
		{SINGLE_QUOTE, `'`},
	})
}

func TestLambdaLexing(t *testing.T) {
	test := newTest(t, "TestLambdaLexing")
	source := `
	(map [1 2 3] \x: (square x))
	`

	test.expect(source, []token{
		{OPEN_PAREN, `(`},
		{IDENTIFIER, `map`},
		{OPEN_BRACKET, `[`},
		{NUMBER, `1`},
		{NUMBER, `2`},
		{NUMBER, `3`},
		{CLOSED_BRACKET, `]`},
		{BACKSLASH, `\`},
		{IDENTIFIER, `x`},
		{COLON, `:`},
		{OPEN_PAREN, `(`},
		{IDENTIFIER, `square`},
		{IDENTIFIER, `x`},
		{CLOSED_PAREN, `)`},
		{CLOSED_PAREN, `)`},
	})
}

func TestTypeDefinitionLexing(t *testing.T) {
	test := newTest(t, "TestTypeDefinitionLexing")
	source := `
	<number -> number -> number>
	add_numbers a b: (add a b)
	`

	test.expect(source, []token{
		{OPEN_ANGLE, `<`},
		{IDENTIFIER, `number`},
		{RIGHT_ARROW, `->`},
		{IDENTIFIER, `number`},
		{RIGHT_ARROW, `->`},
		{IDENTIFIER, `number`},
		{CLOSED_ANGLE, `>`},
		{IDENTIFIER, `add_numbers`},
		{IDENTIFIER, `a`},
		{IDENTIFIER, `b`},
		{COLON, `:`},
		{OPEN_PAREN, `(`},
		{IDENTIFIER, `add`},
		{IDENTIFIER, `a`},
		{IDENTIFIER, `b`},
		{CLOSED_PAREN, `)`},
	})
}

func TestTypeDeclarationLexing(t *testing.T) {
	test := newTest(t, "TestTypeDeclarationLexing")
	source := `
	type name: string

	type person: {
		name: string
	}

	type num_list: [number]
	`

	test.expect(source, []token{
		{TYPE, `type`},
		{IDENTIFIER, `name`},
		{COLON, `:`},
		{IDENTIFIER, `string`},

		{TYPE, `type`},
		{IDENTIFIER, `person`},
		{COLON, `:`},
		{OPEN_BRACE, `{`},
		{IDENTIFIER, `name`},
		{COLON, `:`},
		{IDENTIFIER, `string`},
		{CLOSED_BRACE, `}`},

		{TYPE, `type`},
		{IDENTIFIER, `num_list`},
		{COLON, `:`},
		{OPEN_BRACKET, `[`},
		{IDENTIFIER, `number`},
		{CLOSED_BRACKET, `]`},
	})
}
