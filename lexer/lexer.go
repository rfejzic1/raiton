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
	l.skipWhitespace()

	char, ok := l.current()

	if !ok {
		return l.token(EOF, "")
	}

	if unicode.IsLetter(char) {
		return l.identifierToken()
	} else if unicode.IsDigit(char) {
		return l.numberToken()
	} else if char == '"' {
		return l.stringToken()
	} else {
		return l.specialToken()
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

func (l *lexer) stringToken() Token {
	lexeme := ""

	l.next() // consume quote

	for char, ok := l.current(); ok; char, ok = l.next() {
		if char == '\\' {
			char, ok := l.next()
			if !ok {
				return l.token(EOF, "")
			}

			if char == '"' {
				lexeme += `"`
			} else if char == 'n' {
				lexeme += "\n"
			} else if char == 't' {
				lexeme += "\t"
			} else {
				lexeme += "\\"
				lexeme += string(char)
			}
		} else if char == '"' {
			break
		} else {
			lexeme += string(char)
		}
	}

	if char, ok := l.current(); ok && char == '"' {
		l.next() // consume quote
	}

	return l.longToken(STRING, lexeme)
}

func (l *lexer) specialToken() Token {
	char, _ := l.current()
	lexeme := string(char)
	l.next()

	if char, ok := l.current(); ok {
		extended := lexeme + string(char)
		if tokenType, ok := SYMBOLS[extended]; ok {
			l.next()
			return l.longToken(tokenType, extended)
		}
	}
	if tokenType, ok := SYMBOLS[lexeme]; ok {
		return l.longToken(tokenType, lexeme)
	}

	return l.token(UNKNOWN, lexeme)
}

func (l *lexer) skipWhitespace() {
	for char, ok := l.current(); ok && isWhitespace(char); char, ok = l.next() {
	}
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
		l.consumeComment()
		l.position += 1
		return l.current()
	}
	return 0, false
}

func (l *lexer) current() (rune, bool) {
	if l.ok() {
		l.consumeComment()
		return l.source[l.position], true
	}
	return 0, false
}

func (l *lexer) consumeComment() {
	if isCommentSymbol(l.source[l.position]) {
		for !isLineBreak(l.source[l.position]) {
			l.position += 1
		}
	}
}

func (l *lexer) ok() bool {
	return l.position < len(l.source)
}

func isWhitespace(c rune) bool {
	return isLineBreak(c) || unicode.IsSpace(c)
}

func isLineBreak(c rune) bool {
	return c == '\n'
}

func isCommentSymbol(c rune) bool {
	return c == '#'
}
