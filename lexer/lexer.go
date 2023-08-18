package lexer

import "unicode"

type lexer struct {
	source   []rune
	position int
	line     int
	column   int
}

func New(source []rune) lexer {
	return lexer{
		source:   source,
		line:     1,
		column:   1,
		position: 0,
	}
}

func (l *lexer) Next() Token {
	char, ok := l.current()

	if !ok {
		return l.token(EOF, "")
	}

	if unicode.IsLetter(char) {
		return l.identifierToken()
	} else if unicode.IsDigit(char) {
		return l.numberToken()
	} else {
		return l.token(EOF, "")
	}
}

func (l *lexer) identifierToken() Token {
	literal := ""
	for char, ok := l.current(); ok && (unicode.IsLetter(char) || char == '_'); char, ok = l.next() {
		literal += string(char)
	}
	return l.token(IDENTIFIER, literal)
}

func (l *lexer) numberToken() Token {
	lexeme := ""

	for char, ok := l.current(); ok && unicode.IsDigit(char); char, ok = l.next() {
		lexeme += string(char)
	}

	if char, ok := l.current(); ok && char == '.' {
		lexeme += "."
		l.next()

		for char, ok := l.current(); ok && unicode.IsDigit(char); char, ok = l.next() {
			lexeme += string(char)
		}
	}

	return l.longToken(NUMBER, lexeme)
}

func (l *lexer) token(tokenType TokenType, literal string) Token {
	return Token{
		Literal: literal,
		Type:    tokenType,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *lexer) longToken(tokenType TokenType, literal string) Token {
	return Token{
		Literal: literal,
		Type:    tokenType,
		Line:    l.line,
		Column:  l.column - len(literal),
	}
}

func (l *lexer) next() (rune, bool) {
	if l.ok() {
		l.position += 1
		return l.current()
	}
	return 0, false
}

func (l *lexer) current() (rune, bool) {
	if l.ok() {
		return l.source[l.position], true
	}
	return 0, false
}

func (l *lexer) ok() bool {
	return l.position < len(l.source)
}

func isWhitespace(c rune) bool {
	return !isLineBreak(c) && unicode.IsSpace(c)
}

func isLineBreak(c rune) bool {
	return c == '\n'
}

func isCommentSymbol(c rune) bool {
	return c == '#'
}
