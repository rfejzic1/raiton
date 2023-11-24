package parser

import (
	"fmt"
	"strconv"

	"raiton/ast"
	"raiton/lexer"
	"raiton/token"
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

func (p *Parser) Parse() (ast.Node, error) {
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
	if p.match(token.OPEN_BRACE) {
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
	} else if p.match(token.COLON) {
		scope := &ast.Scope{
			Definitions: make([]*ast.Definition, 0),
			Expressions: make([]ast.Expression, 0),
		}

		p.consume(token.COLON)

		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		scope.Expressions = append(scope.Expressions, expr)

		return scope, nil
	} else {
		return nil, p.unexpected()
	}
}

func (p *Parser) scopeItem(scope *ast.Scope) error {
	if p.match(token.IDENTIFIER) {
		ident := p.identifier()

		switch {
		case p.match(token.DOT):
			selector, err := p.selector(ident)

			if err != nil {
				return err
			}

			scope.Expressions = append(scope.Expressions, selector)
		case p.match(token.COLON):
			definition, err := p.definition(ident)

			if err != nil {
				return err
			}

			scope.Definitions = append(scope.Definitions, definition)
		default:
			scope.Expressions = append(scope.Expressions, ident)
		}
	} else if p.match(token.FUNCTION) {
		funcDef, err := p.functionDefinition()

		if err != nil {
			return err
		}

		scope.Definitions = append(scope.Definitions, funcDef)
	} else {
		expression, err := p.expression()
		if err != nil {
			return err
		}
		scope.Expressions = append(scope.Expressions, expression)
	}

	return nil
}

func (p *Parser) definition(ident *ast.Identifier) (*ast.Definition, error) {
	p.consume(token.COLON)

	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	return &ast.Definition{
		Identifier: *ident,
		Expression: expr,
	}, nil
}

func (p *Parser) functionDefinition() (*ast.Definition, error) {
	if err := p.expect(token.FUNCTION); err != nil {
		return nil, err
	}

	p.consume(token.FUNCTION)

	if err := p.expect(token.IDENTIFIER); err != nil {
		return nil, err
	}

	ident := p.identifier()

	params := []*ast.Identifier{}

	for p.match(token.IDENTIFIER) {
		param := p.identifier()
		params = append(params, param)
	}

	scope, err := p.scope()

	if err != nil {
		return nil, err
	}

	function := &ast.Function{
		Parameters: params,
		Body:       scope,
	}

	return &ast.Definition{
		Identifier: *ident,
		Expression: function,
	}, nil
}

func (p *Parser) expression() (ast.Expression, error) {
	if p.match(token.IDENTIFIER) {
		return p.selector(nil)
	} else if p.match(token.NUMBER) || p.match(token.MINUS) {
		return p.number()
	} else if p.match(token.BOOLEAN) {
		return p.boolean()
	} else if p.match(token.DOUBLE_QUOTE) {
		return p.string()
	} else if p.match(token.SINGLE_QUOTE) {
		return p.string()
	} else if p.match(token.OPEN_BRACKET) {
		return p.arrayOrList()
	} else if p.match(token.OPEN_BRACE) {
		return p.record()
	} else if p.match(token.BACKSLASH) {
		return p.function()
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

func (p *Parser) selector(ident *ast.Identifier) (ast.Expression, error) {
	items := []*ast.SelectorItem{}

	if ident == nil {
		if err := p.expect(token.IDENTIFIER); err != nil {
			return nil, err
		}

		ident = ast.NewIdentifier(p.token.Literal)
		p.consume(token.IDENTIFIER)
	}

	firstItem := ast.NewIdentifierSelector(ident)
	items = append(items, firstItem)

	for p.match(token.DOT) {
		p.consume(token.DOT)

		if p.match(token.IDENTIFIER) {
			ident := p.identifier()
			item := ast.NewIdentifierSelector(ident)
			items = append(items, item)
		} else if p.match(token.NUMBER) {
			num, err := p.unsignedInteger()

			if err != nil {
				return nil, err
			}

			item := ast.NewIndexSelector(num.(*ast.Integer))
			items = append(items, item)
		} else {
			return nil, p.unexpected()
		}
	}

	return &ast.Selector{
		Items: items,
	}, nil
}

func (p *Parser) number() (ast.Expression, error) {
	numberStr := ""

	if p.match(token.MINUS) {
		numberStr += p.token.Literal
		p.consume(token.MINUS)
	}

	if err := p.expect(token.NUMBER); err != nil {
		return nil, err
	}

	numberStr += p.token.Literal
	p.consume(token.NUMBER)

	if p.match(token.DOT) {
		numberStr += p.token.Literal
		p.consume(token.DOT)

		if err := p.expect(token.NUMBER); err != nil {
			return nil, err
		}

		numberStr += p.token.Literal

		value, err := strconv.ParseFloat(numberStr, 64)

		if err != nil {
			return nil, err
		}

		p.consume(token.NUMBER)
		return ast.NewFloat(value), nil
	}

	value, err := strconv.ParseInt(numberStr, 0, 64)

	if err != nil {
		return nil, err
	}

	return ast.NewInteger(value), nil
}

func (p *Parser) unsignedInteger() (ast.Expression, error) {
	if err := p.expect(token.NUMBER); err != nil {
		return nil, err
	}

	value, err := strconv.ParseInt(p.token.Literal, 0, 64)

	if err != nil {
		return nil, err
	}

	p.consume(token.NUMBER)

	return ast.NewInteger(value), nil
}

func (p *Parser) boolean() (ast.Expression, error) {
	value := p.token.Literal
	p.consume(token.BOOLEAN)
	return ast.NewBoolean(value), nil
}

func (p *Parser) string() (ast.Expression, error) {
	quote := p.token.Type

	p.consume(quote)

	if err := p.expect(token.STRING); err != nil {
		return nil, err
	}

	str := ast.NewString(p.token.Literal)

	p.consume(token.STRING)

	if err := p.expect(quote); err != nil {
		return nil, err
	}

	p.consume(quote)

	return str, nil
}

func (p *Parser) arrayOrList() (ast.Expression, error) {
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

		expression, err = p.list()

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

	array := &ast.Array{
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

func (p *Parser) list() (ast.Expression, error) {
	list := &ast.List{
		Elements: []ast.Expression{},
	}

	for !p.match(token.EOF) && !p.match(token.CLOSED_BRACKET) {
		element, err := p.expression()
		if err != nil {
			return nil, err
		}

		list.Elements = append(list.Elements, element)
	}

	return list, nil
}

func (p *Parser) record() (ast.Expression, error) {
	p.consume(token.OPEN_BRACE)

	record := ast.Record{
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

		record.Fields[field] = expression
	}

	if err := p.expect(token.CLOSED_BRACE); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACE)

	return &record, nil
}

func (p *Parser) function() (ast.Expression, error) {
	p.consume(token.BACKSLASH)

	params := []*ast.Identifier{}

	for p.match(token.IDENTIFIER) {
		param := p.identifier()
		params = append(params, param)
	}

	scope, err := p.scope()

	if err != nil {
		return nil, err
	}

	return &ast.Function{
		Parameters: params,
		Body:       scope,
	}, nil
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
