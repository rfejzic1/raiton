package parser

import (
	"fmt"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/token"
)

type Parser struct {
	lex   *lexer.Lexer
	token token.Token
}

func New(lex *lexer.Lexer) Parser {
	return Parser{
		lex: lex,
	}
}

func (p *Parser) Parse() (Expression, error) {
	// The fact that the production method is called
	// means that the current token is matching expecations
	p.consumeAny()
	return p.fileScope()
}

/*** Productions ***/

func (p *Parser) fileScope() (*Scope, error) {
	scope := &Scope{
		definitions:     make([]Definition, 0),
		typeDefinitions: make([]TypeDefinition, 0),
		expressions:     make([]Expression, 0),
	}

	for !p.match(token.EOF) {
		if err := p.scopeItem(scope); err != nil {
			return nil, err
		}
	}

	return scope, nil
}

func (p *Parser) scope() (*Scope, error) {
	scope := &Scope{
		definitions:     make([]Definition, 0),
		typeDefinitions: make([]TypeDefinition, 0),
		expressions:     make([]Expression, 0),
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

func (p *Parser) scopeItem(scope *Scope) error {
	if p.match(token.IDENTIFIER) || p.match(token.OPEN_ANGLE) {
		definition, err := p.definition()
		if err != nil {
			return err
		}
		scope.definitions = append(scope.definitions, definition)
	} else if p.match(token.TYPE) {
		typeDefinition, err := p.typeDefinition()
		if err != nil {
			return err
		}
		scope.typeDefinitions = append(scope.typeDefinitions, typeDefinition)
	} else {
		expression, err := p.expression()
		if err != nil {
			return err
		}
		scope.expressions = append(scope.expressions, expression)
	}

	return nil
}

func (p *Parser) definition() (Definition, error) {
	var err error

	def := Definition{
		parameters: []Identifier{},
	}

	if p.match(token.OPEN_ANGLE) {
		p.consume(token.OPEN_ANGLE)
		def.typeExpression, err = p.typeExpression()
		if err != nil {
			return Definition{}, err
		}
		if err := p.expect(token.CLOSED_ANGLE); err != nil {
			return Definition{}, err
		}
		p.consume(token.CLOSED_ANGLE)
	}

	if err := p.expect(token.IDENTIFIER); err != nil {
		return Definition{}, err
	}

	def.identifier = Identifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	for p.match(token.IDENTIFIER) {
		param := Identifier(p.token.Literal)
		def.parameters = append(def.parameters, param)
		p.consume(token.IDENTIFIER)
	}

	if p.match(token.COLON) {
		p.consume(token.COLON)
		if def.expression, err = p.expression(); err != nil {
			return Definition{}, err
		}
	} else if p.match(token.OPEN_BRACE) {
		if def.expression, err = p.scope(); err != nil {
			return Definition{}, err
		}
	} else {
		return Definition{}, p.unexpected()
	}

	return def, nil
}

func (p *Parser) typeDefinition() (TypeDefinition, error) {
	p.consume(token.TYPE)

	if err := p.expect(token.IDENTIFIER); err != nil {
		return TypeDefinition{}, err
	}

	ident := TypeIdentifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	if err := p.expect(token.COLON); err != nil {
		return TypeDefinition{}, err
	}

	p.consume(token.COLON)

	typeExpression, err := p.typeExpression()

	if err != nil {
		return TypeDefinition{}, err
	}

	return TypeDefinition{
		identifier:     ident,
		typeExpression: typeExpression,
	}, nil
}

func (p *Parser) typeExpression() (TypeExpression, error) {
	var typeExpression TypeExpression

	if p.match(token.IDENTIFIER) {
		typeExpression = TypeIdentifier(p.token.Literal)
		p.consume(token.IDENTIFIER)
	} else if p.match(token.OPEN_PAREN) {
		var err error
		p.consume(token.OPEN_PAREN)
		typeExpression, err = p.typeExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.CLOSED_PAREN); err != nil {
			return nil, err
		}
		p.consume(token.CLOSED_PAREN)
	} else if p.match(token.OPEN_BRACE) {
		p.consume(token.OPEN_BRACE)
		recortType := RecordType{
			fields: map[Identifier]TypeExpression{},
		}

		for p.match(token.IDENTIFIER) {
			field := Identifier(p.token.Literal)
			p.consume(token.IDENTIFIER)
			if err := p.expect(token.COLON); err != nil {
				return nil, err
			}
			typeExpression, err := p.typeExpression()
			if err != nil {
				return nil, err
			}
			recortType.fields[field] = typeExpression
		}

		if err := p.expect(token.CLOSED_BRACE); err != nil {
			return nil, err
		}

		p.consume(token.CLOSED_BRACE)

		typeExpression = recortType
	} else {
		return nil, p.unexpected()
	}

	if p.match(token.RIGHT_ARROW) {
		p.consume(token.RIGHT_ARROW)
		returnTypeExpression, err := p.typeExpression()
		if err != nil {
			return nil, err
		}

		typeExpression = FunctionType{
			parameterType: typeExpression,
			returnType:    returnTypeExpression,
		}
	}

	return typeExpression, nil
}

func (p *Parser) expression() (Expression, error) {
	var expression Expression

	if p.match(token.IDENTIFIER) {
		expression = Identifier(p.token.Literal)
		p.consume(token.IDENTIFIER)
	} else if p.match(token.NUMBER) {
		expression = NumberLiteral(p.token.Literal)
		p.consume(token.NUMBER)
	} else if p.match(token.DOUBLE_QUOTE) {
		p.consume(token.DOUBLE_QUOTE)
		if err := p.expect(token.STRING); err != nil {
			return nil, err
		}
		expression = StringLiteral(p.token.Literal)
		p.consume(token.STRING)
		if err := p.expect(token.DOUBLE_QUOTE); err != nil {
			return nil, err
		}
		p.consume(token.DOUBLE_QUOTE)
	} else if p.match(token.SINGLE_QUOTE) {
		p.consume(token.SINGLE_QUOTE)
		if err := p.expect(token.STRING); err != nil {
			return nil, err
		}
		expression = CharacterLiteral(p.token.Literal)
		p.consume(token.STRING)
		if err := p.expect(token.SINGLE_QUOTE); err != nil {
			return nil, err
		}
		p.consume(token.SINGLE_QUOTE)
	} else if p.match(token.OPEN_BRACKET) {
		p.consume(token.OPEN_BRACKET)

		arrayLiteral := ArrayLiteral{
			elements: []Expression{},
		}

		for !p.match(token.EOF) && !p.match(token.CLOSED_BRACKET) {
			element, err := p.expression()
			if err != nil {
				return nil, err
			}
			arrayLiteral.elements = append(arrayLiteral.elements, element)
		}

		if err := p.expect(token.CLOSED_BRACKET); err != nil {
			return nil, err
		}

		p.consume(token.CLOSED_BRACKET)

		expression = arrayLiteral
	} else if p.match(token.OPEN_BRACE) {
		p.consume(token.OPEN_BRACE)

		recordLiteral := RecordLiteral{
			fields: map[Identifier]Expression{},
		}

		for p.match(token.IDENTIFIER) {
			field := Identifier(p.token.Literal)
			p.consume(token.IDENTIFIER)

			expression, err := p.expression()

			if err != nil {
				return nil, err
			}

			recordLiteral.fields[field] = expression
		}

		if err := p.expect(token.CLOSED_BRACE); err != nil {
			return nil, err
		}

		p.consume(token.CLOSED_BRACE)

		expression = recordLiteral
	} else if p.match(token.BACKSLASH) {
		p.consume(token.BACKSLASH)

		lambdaLiteral := LambdaLiteral{
			parameters: []Identifier{},
		}

		for p.match(token.IDENTIFIER) {
			param := Identifier(p.token.Literal)
			lambdaLiteral.parameters = append(lambdaLiteral.parameters, param)
			p.consume(token.IDENTIFIER)
		}

		if err := p.expect(token.COLON); err != nil {
			return nil, err
		}

		p.consume(token.COLON)

		var err error
		lambdaLiteral.expression, err = p.expression()
		if err != nil {
			return nil, err
		}

		expression = lambdaLiteral
	} else if p.match(token.OPEN_PAREN) {
		p.consume(token.OPEN_PAREN)

		invocation := Invocation{
			arguments: []Expression{},
		}

		for !p.match(token.EOF) && !p.match(token.CLOSED_PAREN) {
			expression, err := p.expression()
			if err != nil {
				return nil, err
			}
			invocation.arguments = append(invocation.arguments, expression)
		}

		if err := p.expect(token.CLOSED_PAREN); err != nil {
			return nil, err
		}

		p.consume(token.CLOSED_PAREN)
		expression = invocation
	} else {
		return nil, p.unexpected()
	}

	return expression, nil
}

/*** Parser utility methods ***/

func (p *Parser) expect(tokenType token.TokenType) error {
	if !p.match(tokenType) {
		return fmt.Errorf("expected %s, but got %s on line %d column %d", tokenType, p.token.Type, p.token.Line, p.token.Column)
	}

	return nil
}

func (p *Parser) unexpected() error {
	return fmt.Errorf("unexpected token `%s` of type %s on line %d column %d", p.token.Literal, p.token.Type, p.token.Line, p.token.Column)
}

func (p *Parser) match(tokenType token.TokenType) bool {
	return p.token.Type == tokenType
}

func (p *Parser) consume(tokenType token.TokenType) {
	if err := p.expect(tokenType); err != nil {
		panic(err)
	}

	p.consumeAny()
}

func (p *Parser) consumeAny() {
	t := p.lex.Next()
	p.token = t
}
