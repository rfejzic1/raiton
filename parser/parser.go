package parser

import (
	"fmt"
	"strconv"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/token"
	"github.com/rfejzic1/raiton/ast"
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
		Definitions:     make([]*ast.Definition, 0),
		TypeDefinitions: make([]*ast.TypeDefinition, 0),
		Expressions:     make([]ast.Expression, 0),
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
		Definitions:     make([]*ast.Definition, 0),
		TypeDefinitions: make([]*ast.TypeDefinition, 0),
		Expressions:     make([]ast.Expression, 0),
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

func (p *Parser) definition() (ast.Definition, error) {
	var err error

	def := ast.Definition{
		Parameters: []*ast.Identifier{},
	}

	if p.match(token.OPEN_ANGLE) {
		p.consume(token.OPEN_ANGLE)
		def.TypeExpression, err = p.typeExpression()
		if err != nil {
			return ast.Definition{}, err
		}
		if err := p.expect(token.CLOSED_ANGLE); err != nil {
			return ast.Definition{}, err
		}
		p.consume(token.CLOSED_ANGLE)
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

func (p *Parser) typeDefinition() (ast.TypeDefinition, error) {
	var err error
	typeDef := ast.TypeDefinition{}

	p.consume(token.TYPE)

	if err := p.expect(token.IDENTIFIER); err != nil {
		return ast.TypeDefinition{}, err
	}

	typeDef.Identifier = ast.TypeIdentifier(p.token.Literal)

	p.consume(token.IDENTIFIER)

	for p.match(token.IDENTIFIER) {
		param := ast.Identifier(p.token.Literal)
		typeDef.Parameters = append(typeDef.Parameters, &param)
		p.consume(token.IDENTIFIER)
	}

	if err := p.expect(token.COLON); err != nil {
		return ast.TypeDefinition{}, err
	}

	p.consume(token.COLON)

	typeDef.TypeExpression, err = p.typeExpression()

	if err != nil {
		return ast.TypeDefinition{}, err
	}

	return typeDef, nil
}

func (p *Parser) typeExpression() (ast.TypeExpression, error) {
	var typeExpression ast.TypeExpression
	var err error

	if p.match(token.IDENTIFIER) {
		typeExpression, err = p.typeIdentifierPath()
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

		typeExpression = &ast.FunctionType{
			ParameterType: typeExpression,
			ReturnType:    returnTypeExpression,
		}
	}

	return typeExpression, nil
}

func (p *Parser) typeIdentifier() *ast.TypeIdentifier {
	ident := ast.TypeIdentifier(p.token.Literal)
	p.consume(token.IDENTIFIER)
	return &ident
}

func (p *Parser) typeIdentifierPath() (ast.TypeExpression, error) {
	identifiers := []*ast.TypeIdentifier{}

	for p.match(token.IDENTIFIER) {
		ident := p.typeIdentifier()
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

	return &ast.TypeIdentifierPath{
		Identifiers: identifiers,
	}, nil
}

func (p *Parser) typeGroup() (ast.TypeExpression, error) {
	p.consume(token.OPEN_PAREN)

	typeExpressions := []ast.TypeExpression{}

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

	return &ast.GroupType{
		TypeExpressions: typeExpressions,
	}, nil
}

func (p *Parser) typeSum() (ast.TypeExpression, error) {
	sumType := ast.SumType{
		Variants: []*ast.SumTypeVariant{},
	}

	for p.match(token.PIPE) {
		p.consume(token.PIPE)

		if err := p.expect(token.IDENTIFIER); err != nil {
			return nil, err
		}

		variant := ast.SumTypeVariant{
			Identifier: ast.Identifier(p.token.Literal),
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

func (p *Parser) typeRecord() (ast.TypeExpression, error) {
	p.consume(token.OPEN_BRACE)
	recortType := ast.RecordType{
		Fields: map[ast.Identifier]ast.TypeExpression{},
	}

	for p.match(token.IDENTIFIER) {
		field := ast.Identifier(p.token.Literal)
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

func (p *Parser) typeArrayOrSlice() (ast.TypeExpression, error) {
	p.consume(token.OPEN_BRACKET)

	var typeExpression ast.TypeExpression

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

		typeExpression = &ast.ArrayType{
			Size:        size,
			ElementType: elementType,
		}
	} else {
		elementType, err := p.typeExpression()

		if err != nil {
			return nil, err
		}

		typeExpression = &ast.SliceType{
			ElementType: elementType,
		}
	}

	if err := p.expect(token.CLOSED_BRACKET); err != nil {
		return nil, err
	}

	p.consume(token.CLOSED_BRACKET)

	return typeExpression, nil
}

func (p *Parser) expression() (ast.Expression, error) {
	if p.match(token.IDENTIFIER) {
		return p.identifierPath()
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

func (p *Parser) identifier() *ast.Identifier {
	ident := ast.Identifier(p.token.Literal)
	p.consume(token.IDENTIFIER)
	return &ident
}

func (p *Parser) identifierPath() (ast.TypeExpression, error) {
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

func (p *Parser) number() ast.Expression {
	num := ast.NewNumberLiteral(p.token.Literal)
	p.consume(token.NUMBER)
	return num
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

	invocation := ast.Invocation{
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
