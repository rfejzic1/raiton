package lexer

import "testing"

func TestIdentifierParsing(t *testing.T) {
	source := `println`

	test := tokenSequenceTest{
		name: "IdentifierParsing",
		expected: []tokenExpect{
			{IDENTIFIER, `println`},
			{EOF, ``},
		},
	}

	test.expect(t, source)
}

func TestNumberParsing(t *testing.T) {
	source := `123`

	test := tokenSequenceTest{
		name: "NumberParsing",
		expected: []tokenExpect{
			{NUMBER, `123`},
			{EOF, ``},
		},
	}

	test.expect(t, source)
}

func TestNumberParsingWithDecimal(t *testing.T) {
	source := `123.4`

	test := tokenSequenceTest{
		name: "NumberParsingWithDecimal",
		expected: []tokenExpect{
			{NUMBER, `123.4`},
			{EOF, ``},
		},
	}

	test.expect(t, source)
}

func TestNumberParsingWithTrailingDecimal(t *testing.T) {
	source := `123.`

	test := tokenSequenceTest{
		name: "NumberParsingWithTrailingDecimal",
		expected: []tokenExpect{
			{NUMBER, "123."},
			{EOF, ``},
		},
	}

	test.expect(t, source)
}

func TestStringParsing(t *testing.T) {
	source := `"Hello, Raiton!"`

	test := tokenSequenceTest{
		name: "StringParsin",
		expected: []tokenExpect{
			{STRING, `Hello, Raiton!`},
			{EOF, ``},
		},
	}

	test.expect(t, source)
}

type tokenSequenceTest struct {
	name     string
	expected []tokenExpect
}

type tokenExpect struct {
	Type    TokenType
	Literal string
}

func (tst *tokenSequenceTest) expect(t *testing.T, source string) {
	src := []rune(source)

	l := New(src)

	for i, et := range tst.expected {
		token := l.Next()

		if token.Type != et.Type {
			t.Fatalf("%s[%d] - wrong token type; expected `%s`, but got `%s` at line %d", tst.name, i, et.Type, token.Type, token.Line)
		}

		if token.Literal != et.Literal {
			t.Fatalf("%s[%d] - wrong token literal; expected `%s`, but got `%s` at line %d", tst.name, i, string(et.Literal), string(token.Literal), token.Line)
		}
	}

}
