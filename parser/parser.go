package parser

import (
	"fmt"
	"strconv"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/token"
)

type Parser struct {
	lex       *lexer.Lexer
	token     token.Token
	peekToken *token.Token
}

func New(lex *lexer.Lexer) Parser {
	return Parser{
		lex: lex,
	}
}

func (p *Parser) Parse() (ast.Expression, error) {
	// The fact that a production method is called
	// means that the current token is matching expecations
	p.nextToken()
	return p.fileScope()
}

/*** Productions ***/

func (p *Parser) fileScope() (*ast.Scope, error) {
	scope := &ast.Scope{
		Definitions: make([]*ast.Definition, 0),
		Expressions: make([]ast.Expression, 0),
	}

	for !p.match(token.EOF) {
		if err := p.scopeItem(scope); err != nil {
			return nil, err
		}
	}

	return scope, nil
}

func (p *Parser) scope() (*ast.Scope, error) {
	scope := &ast.Scope{
		Definitions: make([]*ast.Definition, 0),
		Expressions: make([]ast.Expression, 0),
	}

	p.consume(token.OPEN_BRACE)

	for !p.match(token.EOF) && !p.match(token.CLOSED_BRACE) {
		if err := p.scopeItem(scope); err != nil {
			return nil, err
		}
	}

	if err := p.expect(token.CLOSED_BRACE); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACE)

	return scope, nil
}

func (p *Parser) scopeItem(scope *ast.Scope) error {
	if p.match(token.IDENTIFIER) {
		definition, err := p.definition()
		if err != nil {
			return err
		}
		scope.Definitions = append(scope.Definitions, &definition)
	} else {
		expression, err := p.expression()
		if err != nil {
			return err
		}
		scope.Expressions = append(scope.Expressions, expression)
	}

	return nil
}

func (p *Parser) definition() (ast.Definition, error) {
	var err error

	def := ast.Definition{
		Parameters: []*ast.Identifier{},
	}

	if err := p.expect(token.IDENTIFIER); err != nil {
		return ast.Definition{}, err
	}

	def.Identifier = ast.Identifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	for p.match(token.IDENTIFIER) {
		param := ast.Identifier(p.token.Literal)
		def.Parameters = append(def.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if p.match(token.COLON) {
		p.consume(token.COLON)
		if def.Expression, err = p.expression(); err != nil {
			return ast.Definition{}, err
		}
	} else if p.match(token.OPEN_BRACE) {
		if def.Expression, err = p.scope(); err != nil {
			return ast.Definition{}, err
		}
	} else {
		return ast.Definition{}, p.unexpected()
	}

	return def, nil
}

func (p *Parser) expression() (ast.Expression, error) {
	if p.match(token.IDENTIFIER) {
		return p.identifierPath()
	} else if p.match(token.NUMBER) || p.match(token.MINUS) {
		return p.number()
	} else if p.match(token.BOOLEAN) {
		return p.boolean()
	} else if p.match(token.DOUBLE_QUOTE) {
		return p.string()
	} else if p.match(token.SINGLE_QUOTE) {
		return p.character()
	} else if p.match(token.OPEN_BRACKET) {
		return p.arrayOrSlice()
	} else if p.match(token.OPEN_BRACE) {
		return p.record()
	} else if p.match(token.BACKSLASH) {
		return p.lambda()
	} else if p.match(token.OPEN_PAREN) {
		return p.invocation()
	} else {
		return nil, p.unexpected()
	}
}

func (p *Parser) identifier() *ast.Identifier {
	ident := ast.Identifier(p.token.Literal)
	p.consume(token.IDENTIFIER)
	return &ident
}

func (p *Parser) identifierPath() (ast.Expression, error) {
	identifiers := []*ast.Identifier{}

	for p.match(token.IDENTIFIER) {
		ident := p.identifier()
		identifiers = append(identifiers, ident)

		if p.match(token.DOT) {
			p.consume(token.DOT)

			if err := p.expect(token.IDENTIFIER); err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return &ast.IdentifierPath{
		Identifiers: identifiers,
	}, nil
}

func (p *Parser) number() (ast.Expression, error) {
	numberStr := ""

	if p.match(token.MINUS) {
		numberStr += "-"
		p.consume(token.MINUS)
	}

	if err := p.expect(token.NUMBER); err != nil {
		return nil, err
	}

	numberStr += p.token.Literal

	num := ast.NewNumberLiteral(numberStr)
	p.consume(token.NUMBER)
	return num, nil
}

func (p *Parser) boolean() (ast.Expression, error) {
	// TODO: Implement booleans as sum type instead of having tokens for it
	bool, err := strconv.ParseBool(p.token.Literal)

	if err != nil {
		return nil, err
	}

	p.consume(token.BOOLEAN)

	return ast.NewBooleanLiteral(bool), err
}

func (p *Parser) string() (ast.Expression, error) {
	p.consume(token.DOUBLE_QUOTE)
	if err := p.expect(token.STRING); err != nil {
		return nil, err
	}
	str := ast.NewStringLiteral(p.token.Literal)
	p.consume(token.STRING)
	if err := p.expect(token.DOUBLE_QUOTE); err != nil {
		return nil, err
	}
	p.consume(token.DOUBLE_QUOTE)
	return str, nil
}

func (p *Parser) character() (ast.Expression, error) {
	p.consume(token.SINGLE_QUOTE)
	if err := p.expect(token.STRING); err != nil {
		return nil, err
	}
	char := ast.NewCharacterLiteral(p.token.Literal)
	p.consume(token.STRING)
	if err := p.expect(token.SINGLE_QUOTE); err != nil {
		return nil, err
	}
	p.consume(token.SINGLE_QUOTE)
	return char, nil
}

func (p *Parser) arrayOrSlice() (ast.Expression, error) {
	var expression ast.Expression

	p.consume(token.OPEN_BRACKET)

	if p.match(token.NUMBER) {
		p.peek()

		if p.peekMatch(token.COLON) {
			var err error

			expression, err = p.array()

			if err != nil {
				return nil, err
			}
		}
	}

	if expression == nil {
		var err error

		expression, err = p.slice()

		if err != nil {
			return nil, err
		}
	}

	if err := p.expect(token.CLOSED_BRACKET); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACKET)

	return expression, nil
}

func (p *Parser) array() (ast.Expression, error) {
	size, err := parseArraySize(p.token.Literal)

	if err != nil {
		return nil, err
	}

	p.consume(token.NUMBER)

	if err := p.expect(token.COLON); err != nil {
		return nil, err
	}

	p.consume(token.COLON)

	array := &ast.ArrayLiteral{
		Size:     size,
		Elements: []ast.Expression{},
	}

	for !p.match(token.EOF) && !p.match(token.CLOSED_BRACKET) {
		element, err := p.expression()
		if err != nil {
			return nil, err
		}

		array.Elements = append(array.Elements, element)
	}

	return array, nil
}

func (p *Parser) slice() (ast.Expression, error) {
	slice := &ast.SliceLiteral{
		Elements: []ast.Expression{},
	}

	for !p.match(token.EOF) && !p.match(token.CLOSED_BRACKET) {
		element, err := p.expression()
		if err != nil {
			return nil, err
		}

		slice.Elements = append(slice.Elements, element)
	}

	return slice, nil
}

func (p *Parser) record() (ast.Expression, error) {
	p.consume(token.OPEN_BRACE)

	recordLiteral := ast.RecordLiteral{
		Fields: map[ast.Identifier]ast.Expression{},
	}

	for p.match(token.IDENTIFIER) {
		field := ast.Identifier(p.token.Literal)
		p.consume(token.IDENTIFIER)

		if err := p.expect(token.COLON); err != nil {
			return nil, err
		}

		p.consume(token.COLON)

		expression, err := p.expression()

		if err != nil {
			return nil, err
		}

		recordLiteral.Fields[field] = expression
	}

	if err := p.expect(token.CLOSED_BRACE); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACE)

	return &recordLiteral, nil
}

func (p *Parser) lambda() (ast.Expression, error) {
	p.consume(token.BACKSLASH)

	lambdaLiteral := ast.LambdaLiteral{
		Parameters: []*ast.Identifier{},
	}

	var err error

	for p.match(token.IDENTIFIER) {
		param := ast.Identifier(p.token.Literal)
		lambdaLiteral.Parameters = append(lambdaLiteral.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if p.match(token.COLON) {
		p.consume(token.COLON)
		if lambdaLiteral.Expression, err = p.expression(); err != nil {
			return &ast.Definition{}, err
		}
	} else if p.match(token.OPEN_BRACE) {
		if lambdaLiteral.Expression, err = p.scope(); err != nil {
			return &ast.Definition{}, err
		}
	} else {
		return &ast.Definition{}, p.unexpected()
	}

	return &lambdaLiteral, nil
}

func (p *Parser) invocation() (ast.Expression, error) {
	p.consume(token.OPEN_PAREN)

	invocation := ast.Application{
		Arguments: []ast.Expression{},
	}

	for !p.match(token.EOF) && !p.match(token.CLOSED_PAREN) {
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		invocation.Arguments = append(invocation.Arguments, expression)
	}

	if err := p.expect(token.CLOSED_PAREN); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_PAREN)

	return &invocation, nil
}

/*** Parser utility methods ***/

func parseArraySize(literal string) (uint64, error) {
	size, err := strconv.ParseUint(literal, 10, 0)

	if err != nil {
		return 0, fmt.Errorf("expected a non-negative integer, but got %s", literal)
	}

	return size, nil
}

func (p *Parser) unexpected() error {
	return fmt.Errorf("unexpected token `%s` of type %s on line %d column %d", p.token.Literal, p.token.Type, p.token.Line, p.token.Column)
}

func (p *Parser) consume(tokenType token.TokenType) {
	if err := p.expect(tokenType); err != nil {
		panic(err)
	}

	if p.peekToken != nil {
		p.token = *p.peekToken
		p.peekToken = nil
	} else {
		p.nextToken()
	}
}

func (p *Parser) match(tokenType token.TokenType) bool {
	return p.token.Type == tokenType
}

func (p *Parser) expect(tokenType token.TokenType) error {
	if !p.match(tokenType) {
		return fmt.Errorf("expected %s, but got %s on line %d column %d", tokenType, p.token.Type, p.token.Line, p.token.Column)
	}

	return nil
}

func (p *Parser) peek() {
	if p.peekToken == nil {
		t := p.lex.Next()
		p.peekToken = &t
	}
}

func (p *Parser) peekMatch(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) nextToken() {
	t := p.lex.Next()
	p.token = t
}
