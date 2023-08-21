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
	p.consume()
	return p.fileScope()
}

/*** Productions ***/

func (p *Parser) fileScope() (*Scope, error) {
	return p.scopeContent()
}

func (p *Parser) scope() (*Scope, error) {
	p.consume() // consume OPEN_BRACE

	s, err := p.scopeContent()

	// expect '}'
	if err := p.expect(token.CLOSED_BRACE); err != nil {
		return nil, err
	}

	p.consume() // consume CLOSED_BRACE

	return s, err
}

func (p *Parser) scopeContent() (*Scope, error) {
	scope := &Scope{
		definitions:     make([]Definition, 0),
		typeDefinitions: make([]TypeDefinition, 0),
		expressions:     make([]Expression, 0),
	}

	for !p.match(token.EOF) {
		if p.match(token.IDENTIFIER) || p.match(token.OPEN_ANGLE) {
			definition, err := p.definition()
			if err != nil {
				return nil, err
			}
			scope.definitions = append(scope.definitions, definition)
		} else if p.match(token.TYPE) {
			typeDefinition, err := p.typeDefinition()
			if err != nil {
				return nil, err
			}
			scope.typeDefinitions = append(scope.typeDefinitions, typeDefinition)
		} else {
			expression, err := p.expression()
			if err != nil {
				return nil, err
			}
			scope.expressions = append(scope.expressions, expression)
		}
	}

	return scope, nil
}

func (p *Parser) definition() (Definition, error) {
	var err error

	def := Definition{
		parameters: []Identifier{},
	}

	if p.match(token.OPEN_ANGLE) {
		p.consume() // consume OPEN_ANGlE
		def.typeExpression, err = p.typeExpression()
		if err != nil {
			return Definition{}, err
		}
		if err := p.expect(token.CLOSED_ANGLE); err != nil {
			return Definition{}, err
		}
		p.consume() // consume CLOSED_ANGLE
	}

	if err := p.expect(token.IDENTIFIER); err != nil {
		return Definition{}, err
	}

	def.identifier = Identifier(p.token.Literal)

	p.consume()

	for p.match(token.IDENTIFIER) {
		param := Identifier(p.token.Literal)
		def.parameters = append(def.parameters, param)
		p.consume()
	}

	if p.match(token.COLON) {
		p.consume() // consume COLON
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
	p.consume() // consume TYPE

	if err := p.expect(token.IDENTIFIER); err != nil {
		return TypeDefinition{}, err
	}

	ident := TypeIdentifier(p.token.Literal)

	p.consume() // consume IDENTIFIER

	if err := p.expect(token.COLON); err != nil {
		return TypeDefinition{}, err
	}

	p.consume() // consume COLON

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
	var err error

	if p.match(token.IDENTIFIER) {
		typeExpression = TypeIdentifier(p.token.Literal)
		p.consume() // consume IDENTIFIER
	} else if p.match(token.OPEN_PAREN) {
		p.consume() // consume OPEN_PAREN
		typeExpression, err = p.typeExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.CLOSED_PAREN); err != nil {
			return nil, err
		}
		p.consume() // consume CLOSED_PAREN
	} else if p.match(token.OPEN_BRACE) {
		p.consume() // consume OPEN_BRACE
		recortType := RecordType{
			fields: map[Identifier]TypeExpression{},
		}

		for p.match(token.IDENTIFIER) {
			field := Identifier(p.token.Literal)
			p.consume()
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

		p.consume() // consume CLOSED_BRACE
	} else {
		return nil, p.unexpected()
	}

	if p.match(token.RIGHT_ARROW) {
		p.consume() // consume RIGHT_ARROW
		returnTypeExpression, err := p.typeExpression()
		if err != nil {
			return nil, err
		}

		return FunctionType{
			parameterType: typeExpression,
			returnType:    returnTypeExpression,
		}, nil
	}

	return typeExpression, nil
}

func (p *Parser) expression() (Expression, error) {
	// parse expressions
	return nil, nil
}

/*** Parser utility methods ***/

func (p *Parser) expectNext(tokenType token.TokenType) error {
	p.consume()
	return p.expect(tokenType)
}

func (p *Parser) expect(tokenType token.TokenType) error {
	if p.match(tokenType) {
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

func (p *Parser) consume() token.Token {
	t := p.lex.Next()
	p.token = t
	return p.token
}
