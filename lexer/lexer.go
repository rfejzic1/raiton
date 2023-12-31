package lexer

import "raiton/token"

type lexMode uint

const (
	NORMAL_MODE lexMode = iota
	SEQUENCE_MODE
)

type Lexer struct {
	source   string
	position int
	line     int
	column   int
	mode     lexMode
	modeChar byte
}

func New(source string) Lexer {
	return Lexer{
		source:   source,
		line:     1,
		column:   1,
		position: 0,
		mode:     NORMAL_MODE,
		modeChar: 0,
	}
}

func (l *Lexer) Next() token.Token {
	switch l.mode {
	case NORMAL_MODE:
		return l.normalMode()
	case SEQUENCE_MODE:
		return l.sequenceMode()
	default:
		return l.token(token.ILLEGAL, "")
	}
}

func (l *Lexer) normalMode() token.Token {
	l.skipWhitespace()

	char, ok := l.current()

	if !ok {
		return l.token(token.EOF, "")
	}

	if isUnderscore(char) || isAlpha(char) {
		return l.identifierToken()
	} else if isDigit(char) {
		return l.numberToken()
	} else if isQuote(char) {
		token := l.specialToken()
		l.mode = SEQUENCE_MODE
		l.modeChar = char
		return token
	} else {
		return l.specialToken()
	}
}

func (l *Lexer) sequenceMode() token.Token {
	char, ok := l.current()

	if !ok {
		return l.token(token.EOF, "")
	}

	if char == l.modeChar {
		token := l.specialToken()
		l.mode = NORMAL_MODE
		return token
	}

	return l.stringToken()
}

func (l *Lexer) identifierToken() token.Token {
	literal := ""

	if char, ok := l.current(); ok && (isAlpha(char) || isUnderscore(char)) {
		literal += string(char)
		l.next()
	}

	for char, ok := l.current(); ok && (isAlpha(char) || isUnderscore(char) || isDigit(char)); char, ok = l.next() {
		literal += string(char)
	}

	if char, ok := l.current(); ok && char == '!' {
		literal += string(char)
		l.next()
	}

	tokenType, ok := token.KEYWORDS[literal]
	if !ok {
		tokenType = token.IDENTIFIER
	}

	return l.longToken(tokenType, literal)
}

func (l *Lexer) numberToken() token.Token {
	lexeme := ""

	for char, ok := l.current(); ok && isDigit(char); char, ok = l.next() {
		lexeme += string(char)
	}

	return l.longToken(token.NUMBER, lexeme)
}

func (l *Lexer) stringToken() token.Token {
	lexeme := ""

	for char, ok := l.current(); ok; char, ok = l.next() {
		if char == '\\' {
			char, ok := l.next()
			if !ok {
				return l.token(token.EOF, "")
			}

			if char == '"' {
				lexeme += `"`
			} else if char == '\'' {
				lexeme += `'`
			} else if char == 'n' {
				lexeme += "\n"
			} else if char == 't' {
				lexeme += "\t"
			} else {
				lexeme += "\\"
				lexeme += string(char)
			}
		} else if char == l.modeChar {
			break
		} else {
			lexeme += string(char)
		}
	}

	return l.longToken(token.STRING, lexeme)
}

func (l *Lexer) specialToken() token.Token {
	char, _ := l.current()
	lexeme := string(char)
	l.next()

	if char, ok := l.current(); ok {
		extended := lexeme + string(char)
		if tokenType, ok := token.SYMBOLS[extended]; ok {
			l.next()
			return l.longToken(tokenType, extended)
		}
	}

	if tokenType, ok := token.SYMBOLS[lexeme]; ok {
		return l.longToken(tokenType, lexeme)
	}

	return l.longToken(token.ILLEGAL, lexeme)
}

func (l *Lexer) skipWhitespace() {
	for char, ok := l.current(); ok && isWhitespace(char); char, ok = l.next() {
		if isLineBreak(char) {
			l.line += 1
			l.column = 0
		}
	}
}

func (l *Lexer) token(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Literal: literal,
		Type:    tokenType,
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) longToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Literal: literal,
		Type:    tokenType,
		Line:    l.line,
		Column:  l.column - len(literal),
	}
}

func (l *Lexer) next() (byte, bool) {
	if l.ok() {
		l.consumeComment()
		l.position += 1
		l.column += 1
		return l.current()
	}
	return 0, false
}

func (l *Lexer) current() (byte, bool) {
	if l.ok() {
		l.consumeComment()
		return l.source[l.position], true
	}
	return 0, false
}

func (l *Lexer) consumeComment() {
	if isCommentSymbol(l.source[l.position]) {
		for !isLineBreak(l.source[l.position]) {
			l.position += 1
		}
	}
}

func (l *Lexer) ok() bool {
	return l.position < len(l.source)
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isDash(c byte) bool {
	return c == '-'
}

func isUnderscore(c byte) bool {
	return c == '_'
}

func isQuote(c byte) bool {
	return c == '"' || c == '\''
}

func isWhitespace(c byte) bool {
	return isLineBreak(c) || c == ' ' || c == '\t'
}

func isLineBreak(c byte) bool {
	return c == '\n'
}

func isCommentSymbol(c byte) bool {
	return c == '#'
}
