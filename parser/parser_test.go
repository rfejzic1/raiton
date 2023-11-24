package parser

import (
	"testing"

	"raiton/ast"
	"raiton/lexer"
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
			ast.NewString("this is a string"),
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

	# number expression; negative float
	-3.14
	`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewInteger(5),
			ast.NewFloat(2.65),
			ast.NewInteger(-1),
			ast.NewFloat(-3.14),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionArray(t *testing.T) {
	source := `[3: 1 2 3]`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewArray(
				3,
				ast.NewInteger(1),
				ast.NewInteger(2),
				ast.NewInteger(3),
			),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionList(t *testing.T) {
	source := `[1 2 3]`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewList(
				ast.NewInteger(1),
				ast.NewInteger(2),
				ast.NewInteger(3),
			),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionInvocation(t *testing.T) {
	source := `(println "Hello, World")`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewApplication(
				ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("println"))),
				ast.NewString("Hello, World"),
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
				Expression: ast.NewString("Tojuro"),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestFunctionDefinitionWithSingleExpression(t *testing.T) {
	source := `
	fn add_two x: (add x 2)
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				Identifier: ast.Identifier("add_two"),
				Expression: &ast.Function{
					Parameters: []*ast.Identifier{
						ast.NewIdentifier("x"),
					},
					Body: ast.ScopeExpressions(
						ast.NewApplication(
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("add"))),
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("x"))),
							ast.NewInteger(2),
						),
					),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestFunctionDefinitionWithScope(t *testing.T) {
	source := `
	fn add_three x { (add x 3) }
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				Identifier: ast.Identifier("add_three"),
				Expression: &ast.Function{
					Parameters: []*ast.Identifier{
						ast.NewIdentifier("x"),
					},
					Body: ast.ScopeExpressions(
						ast.NewApplication(
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("add"))),
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("x"))),
							ast.NewInteger(3),
						),
					),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}
