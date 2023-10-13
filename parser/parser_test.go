package parser

import (
	"testing"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/lexer"
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

func TestDefinitionWithSingleExpression(t *testing.T) {
	source := `
	name: "Tojuro"
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				Identifier: ast.Identifier("name"),
				Expression: ast.NewStringLiteral("Tojuro"),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestDefinitionWithScope(t *testing.T) {
	source := `
	age { 24 }
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				Identifier: ast.Identifier("age"),
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

func TestFunctionDefinitionWithSingleExpression(t *testing.T) {
	source := `
	add_two x: (add x 2)
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
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

func TestFunctionDefinitionWithScope(t *testing.T) {
	source := `
	add_three x { (add x 3) }
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
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
