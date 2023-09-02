package parser

import (
	"testing"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/ast"
)

func parseAndCompare(t *testing.T, source string, expected ast.Expression) {
	l := lexer.New(source)
	p := New(&l)
	got, err := p.Parse()

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	comp := ast.NewComparator(got)

	if err := comp.Compare(expected); err != nil {
		t.Fatalf("assertion failed: %s", err)
	}
}

func TestExpressionString(t *testing.T) {
	source := `"this is a string"`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewStringLiteral("this is a string"),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionCharacter(t *testing.T) {
	source := `'c'`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewCharacterLiteral("c"),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionNumber(t *testing.T) {
	source := `
	# number expression; positive integer
	5

	# number expression; positive float
	2.65

	# number expression; negative integer
	-1
	`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewNumberLiteral("5"),
			ast.NewNumberLiteral("2.65"),
			ast.NewNumberLiteral("-1"),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionArray(t *testing.T) {
	source := `[3: 1 2 3]`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewArrayLiteral(
				3,
				ast.NewNumberLiteral("1"),
				ast.NewNumberLiteral("2"),
				ast.NewNumberLiteral("3"),
			),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionSlice(t *testing.T) {
	source := `[1 2 3]`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewSliceLiteral(
				ast.NewNumberLiteral("1"),
				ast.NewNumberLiteral("2"),
				ast.NewNumberLiteral("3"),
			),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionInvocation(t *testing.T) {
	source := `(println "Hello, World")`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewInvocation(
				ast.NewIdentifierPath(ast.NewIdentifier("println")),
				ast.NewStringLiteral("Hello, World"),
			),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestDefinitionTypedWithSingleExpression(t *testing.T) {
	source := `
	<string>
	name: "Tojuro"
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				TypeExpression: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("string")),
				Identifier:     ast.Identifier("name"),
				Expression:     ast.NewStringLiteral("Tojuro"),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestDefinitionTypedWithScope(t *testing.T) {
	source := `
	<number>
	age { 24 }
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				TypeExpression: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
				Identifier:     ast.Identifier("age"),
				Expression: &ast.Scope{
					Expressions: []ast.Expression{
						ast.NewNumberLiteral("24"),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestFunctionDefinitionTypedWithSingleExpression(t *testing.T) {
	source := `
	<number -> number>
	add_two x: (add x 2)
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				TypeExpression: &ast.FunctionType{
					ParameterType: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
					ReturnType:    ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
				},
				Identifier: ast.Identifier("add_two"),
				Parameters: []*ast.Identifier{
					ast.NewIdentifier("x"),
				},
				Expression: ast.NewInvocation(
					ast.NewIdentifierPath(ast.NewIdentifier("add")),
					ast.NewIdentifierPath(ast.NewIdentifier("x")),
					ast.NewNumberLiteral("2"),
				),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestFunctionDefinitionTypedWithScope(t *testing.T) {
	source := `
	<number -> number>
	add_three x { (add x 3) }
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				TypeExpression: &ast.FunctionType{
					ParameterType: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
					ReturnType:    ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
				},
				Identifier: ast.Identifier("add_three"),
				Parameters: []*ast.Identifier{
					ast.NewIdentifier("x"),
				},
				Expression: &ast.Scope{
					Expressions: []ast.Expression{
						ast.NewInvocation(
							ast.NewIdentifierPath(ast.NewIdentifier("add")),
							ast.NewIdentifierPath(ast.NewIdentifier("x")),
							ast.NewNumberLiteral("3"),
						),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionAlias(t *testing.T) {
	source := `
	type name: string
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier:     ast.TypeIdentifier("name"),
				TypeExpression: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("string")),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionArray(t *testing.T) {
	source := `
	type numArray: [3: number]
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("numArray"),
				TypeExpression: &ast.ArrayType{
					Size:        3,
					ElementType: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionSlice(t *testing.T) {
	source := `
	type numSlice: [number]
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("numSlice"),
				TypeExpression: &ast.SliceType{
					ElementType: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionRecord(t *testing.T) {
	source := `
	type person: {
		name: string
		age: number
	}
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("person"),
				TypeExpression: &ast.RecordType{
					Fields: map[ast.Identifier]ast.TypeExpression{
						ast.Identifier("name"): ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("string")),
						ast.Identifier("age"):  ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionSum(t *testing.T) {
	source := `
	type color: | Red | Green | Blue | RGB: { r:number g:number b:number }
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("color"),
				TypeExpression: &ast.SumType{
					Variants: []*ast.SumTypeVariant{
						{
							Identifier: ast.Identifier("Red"),
						},
						{
							Identifier: ast.Identifier("Green"),
						},
						{
							Identifier: ast.Identifier("Blue"),
						},
						{
							Identifier: ast.Identifier("RGB"),
							TypeExpression: &ast.RecordType{
								Fields: map[ast.Identifier]ast.TypeExpression{
									ast.Identifier("r"): ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
									ast.Identifier("g"): ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
									ast.Identifier("b"): ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("number")),
								},
							},
						},
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionParametrized(t *testing.T) {
	source := `
	type option T: | Some: T | None
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("option"),
				Parameters: []*ast.Identifier{
					ast.NewIdentifier("T"),
				},
				TypeExpression: &ast.SumType{
					Variants: []*ast.SumTypeVariant{
						{
							Identifier:     ast.Identifier("Some"),
							TypeExpression: ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("T")),
						},
						{
							Identifier: ast.Identifier("None"),
						},
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionGroup(t *testing.T) {
	source := `
	type stringOption: (option string)
	`

	expected := ast.Scope{
		TypeDefinitions: []*ast.TypeDefinition{
			{
				Identifier: ast.TypeIdentifier("stringOption"),
				TypeExpression: &ast.GroupType{
					TypeExpressions: []ast.TypeExpression{
						ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("option")),
						ast.NewTypeIdentifierPath(ast.NewTypeIdentifier("string")),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}
