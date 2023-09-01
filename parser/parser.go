package parser

import (
	"fmt"
	"strconv"

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

func (p *Parser) Parse() (Expression, error) {
	// The fact that a production method is called
	// means that the current token is matching expecations
	p.nextToken()
	return p.fileScope()
}

/*** Productions ***/

func (p *Parser) fileScope() (*Scope, error) {
	scope := &Scope{
		Definitions:     make([]*Definition, 0),
		TypeDefinitions: make([]*TypeDefinition, 0),
		Expressions:     make([]Expression, 0),
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
		Definitions:     make([]*Definition, 0),
		TypeDefinitions: make([]*TypeDefinition, 0),
		Expressions:     make([]Expression, 0),
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
		scope.Definitions = append(scope.Definitions, &definition)
	} else if p.match(token.TYPE) {
		typeDefinition, err := p.typeDefinition()
		if err != nil {
			return err
		}
		scope.TypeDefinitions = append(scope.TypeDefinitions, &typeDefinition)
	} else {
		expression, err := p.expression()
		if err != nil {
			return err
		}
		scope.Expressions = append(scope.Expressions, expression)
	}

	return nil
}

func (p *Parser) definition() (Definition, error) {
	var err error

	def := Definition{
		Parameters: []*Identifier{},
	}

	if p.match(token.OPEN_ANGLE) {
		p.consume(token.OPEN_ANGLE)
		def.TypeExpression, err = p.typeExpression()
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

	def.Identifier = Identifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	for p.match(token.IDENTIFIER) {
		param := Identifier(p.token.Literal)
		def.Parameters = append(def.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if p.match(token.COLON) {
		p.consume(token.COLON)
		if def.Expression, err = p.expression(); err != nil {
			return Definition{}, err
		}
	} else if p.match(token.OPEN_BRACE) {
		if def.Expression, err = p.scope(); err != nil {
			return Definition{}, err
		}
	} else {
		return Definition{}, p.unexpected()
	}

	return def, nil
}

func (p *Parser) typeDefinition() (TypeDefinition, error) {
	var err error
	typeDef := TypeDefinition{}

	p.consume(token.TYPE)

	if err := p.expect(token.IDENTIFIER); err != nil {
		return TypeDefinition{}, err
	}

	typeDef.Identifier = TypeIdentifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	for p.match(token.IDENTIFIER) {
		param := Identifier(p.token.Literal)
		typeDef.Parameters = append(typeDef.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if err := p.expect(token.COLON); err != nil {
		return TypeDefinition{}, err
	}

	p.consume(token.COLON)

	typeDef.TypeExpression, err = p.typeExpression()

	if err != nil {
		return TypeDefinition{}, err
	}

	return typeDef, nil
}

func (p *Parser) typeExpression() (TypeExpression, error) {
	var typeExpression TypeExpression
	var err error

	if p.match(token.IDENTIFIER) {
		typeExpression = p.typeIdentifier()
	} else if p.match(token.PIPE) {
		typeExpression, err = p.typeSum()
	} else if p.match(token.OPEN_PAREN) {
		typeExpression, err = p.typeGroup()
	} else if p.match(token.OPEN_BRACKET) {
		typeExpression, err = p.typeArrayOrSlice()
	} else if p.match(token.OPEN_BRACE) {
		typeExpression, err = p.typeRecord()
	} else {
		return nil, p.unexpected()
	}

	if err != nil {
		return nil, err
	}

	if p.match(token.RIGHT_ARROW) {
		p.consume(token.RIGHT_ARROW)
		returnTypeExpression, err := p.typeExpression()
		if err != nil {
			return nil, err
		}

		typeExpression = &FunctionType{
			ParameterType: typeExpression,
			ReturnType:    returnTypeExpression,
		}
	}

	return typeExpression, nil
}

func (p *Parser) typeIdentifier() TypeExpression {
	ident := TypeIdentifier(p.token.Literal)
	p.consume(token.IDENTIFIER)
	return &ident
}

func (p *Parser) typeGroup() (TypeExpression, error) {
	p.consume(token.OPEN_PAREN)

	typeExpressions := []TypeExpression{}

	for !p.match(token.EOF) && !p.match(token.CLOSED_PAREN) {
		typeExpression, err := p.typeExpression()
		if err != nil {
			return nil, err
		}

		typeExpressions = append(typeExpressions, typeExpression)
	}

	if err := p.expect(token.CLOSED_PAREN); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_PAREN)

	return &GroupType{
		TypeExpressions: typeExpressions,
	}, nil
}

func (p *Parser) typeSum() (TypeExpression, error) {
	sumType := SumType{
		Variants: []*SumTypeVariant{},
	}

	for p.match(token.PIPE) {
		p.consume(token.PIPE)

		if err := p.expect(token.IDENTIFIER); err != nil {
			return nil, err
		}

		variant := SumTypeVariant{
			Identifier: Identifier(p.token.Literal),
		}

		p.consume(token.IDENTIFIER)

		if p.match(token.COLON) {
			p.consume(token.COLON)

			typeExpression, err := p.typeExpression()
			if err != nil {
				return nil, err
			}

			variant.TypeExpression = typeExpression
		}

		sumType.Variants = append(sumType.Variants, &variant)
	}

	return &sumType, nil
}

func (p *Parser) typeRecord() (TypeExpression, error) {
	p.consume(token.OPEN_BRACE)
	recortType := RecordType{
		Fields: map[Identifier]TypeExpression{},
	}

	for p.match(token.IDENTIFIER) {
		field := Identifier(p.token.Literal)
		p.consume(token.IDENTIFIER)
		if err := p.expect(token.COLON); err != nil {
			return nil, err
		}
		p.consume(token.COLON)
		typeExpression, err := p.typeExpression()
		if err != nil {
			return nil, err
		}
		recortType.Fields[field] = typeExpression
	}

	if err := p.expect(token.CLOSED_BRACE); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACE)

	return &recortType, nil
}

func (p *Parser) typeArrayOrSlice() (TypeExpression, error) {
	p.consume(token.OPEN_BRACKET)

	var typeExpression TypeExpression

	if p.match(token.NUMBER) {
		size, err := parseArraySize(p.token.Literal)

		if err != nil {
			return nil, err
		}

		p.consume(token.NUMBER)

		if err := p.expect(token.COLON); err != nil {
			return nil, err
		}

		p.consume(token.COLON)

		elementType, err := p.typeExpression()

		if err != nil {
			return nil, err
		}

		typeExpression = &ArrayType{
			Size:        size,
			ElementType: elementType,
		}
	} else {
		elementType, err := p.typeExpression()

		if err != nil {
			return nil, err
		}

		typeExpression = &SliceType{
			ElementType: elementType,
		}
	}

	if err := p.expect(token.CLOSED_BRACKET); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACKET)

	return typeExpression, nil
}

func (p *Parser) expression() (Expression, error) {
	if p.match(token.IDENTIFIER) {
		return p.identifier(), nil
	} else if p.match(token.NUMBER) {
		return p.number(), nil
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

func (p *Parser) identifier() Expression {
	ident := Identifier(p.token.Literal)
	p.consume(token.IDENTIFIER)
	return &ident
}

func (p *Parser) number() Expression {
	num := NewNumberLiteral(p.token.Literal)
	p.consume(token.NUMBER)
	return num
}

func (p *Parser) string() (Expression, error) {
	p.consume(token.DOUBLE_QUOTE)
	if err := p.expect(token.STRING); err != nil {
		return nil, err
	}
	str := NewStringLiteral(p.token.Literal)
	p.consume(token.STRING)
	if err := p.expect(token.DOUBLE_QUOTE); err != nil {
		return nil, err
	}
	p.consume(token.DOUBLE_QUOTE)
	return str, nil
}

func (p *Parser) character() (Expression, error) {
	p.consume(token.SINGLE_QUOTE)
	if err := p.expect(token.STRING); err != nil {
		return nil, err
	}
	char := NewCharacterLiteral(p.token.Literal)
	p.consume(token.STRING)
	if err := p.expect(token.SINGLE_QUOTE); err != nil {
		return nil, err
	}
	p.consume(token.SINGLE_QUOTE)
	return char, nil
}

func (p *Parser) arrayOrSlice() (Expression, error) {
	var expression Expression

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

func (p *Parser) array() (Expression, error) {
	size, err := parseArraySize(p.token.Literal)

	if err != nil {
		return nil, err
	}

	p.consume(token.NUMBER)

	if err := p.expect(token.COLON); err != nil {
		return nil, err
	}

	p.consume(token.COLON)

	array := &ArrayLiteral{
		Size:     size,
		Elements: []Expression{},
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

func (p *Parser) slice() (Expression, error) {
	slice := &SliceLiteral{
		Elements: []Expression{},
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

func (p *Parser) record() (Expression, error) {
	p.consume(token.OPEN_BRACE)

	recordLiteral := RecordLiteral{
		Fields: map[Identifier]Expression{},
	}

	for p.match(token.IDENTIFIER) {
		field := Identifier(p.token.Literal)
		p.consume(token.IDENTIFIER)

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

func (p *Parser) lambda() (Expression, error) {
	p.consume(token.BACKSLASH)

	lambdaLiteral := LambdaLiteral{
		Parameters: []*Identifier{},
	}

	var err error

	for p.match(token.IDENTIFIER) {
		param := Identifier(p.token.Literal)
		lambdaLiteral.Parameters = append(lambdaLiteral.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if p.match(token.COLON) {
		p.consume(token.COLON)
		if lambdaLiteral.Expression, err = p.expression(); err != nil {
			return &Definition{}, err
		}
	} else if p.match(token.OPEN_BRACE) {
		if lambdaLiteral.Expression, err = p.scope(); err != nil {
			return &Definition{}, err
		}
	} else {
		return &Definition{}, p.unexpected()
	}

	return &lambdaLiteral, nil
}

func (p *Parser) invocation() (Expression, error) {
	p.consume(token.OPEN_PAREN)

	invocation := Invocation{
		Arguments: []Expression{},
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
