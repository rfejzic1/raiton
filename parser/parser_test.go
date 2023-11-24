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

	# number expression; negative float
	-3.14
	`

	expected := ast.Scope{
		Expressions: []ast.Expression{
			ast.NewIntegerLiteral(5),
			ast.NewFloatLiteral(2.65),
			ast.NewIntegerLiteral(-1),
			ast.NewFloatLiteral(-3.14),
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
				ast.NewIntegerLiteral(1),
				ast.NewIntegerLiteral(2),
				ast.NewIntegerLiteral(3),
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
				ast.NewIntegerLiteral(1),
				ast.NewIntegerLiteral(2),
				ast.NewIntegerLiteral(3),
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

func TestFunctionDefinitionWithSingleExpression(t *testing.T) {
	source := `
	fn add_two x: (add x 2)
	`

	expected := ast.Scope{
		Definitions: []*ast.Definition{
			{
				Identifier: ast.Identifier("add_two"),
				Expression: &ast.FunctionLiteral{
					Parameters: []*ast.Identifier{
						ast.NewIdentifier("x"),
					},
					Body: ast.ScopeExpressions(
						ast.NewApplication(
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("add"))),
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("x"))),
							ast.NewIntegerLiteral(2),
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
				Expression: &ast.FunctionLiteral{
					Parameters: []*ast.Identifier{
						ast.NewIdentifier("x"),
					},
					Body: ast.ScopeExpressions(
						ast.NewApplication(
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("add"))),
							ast.NewSelector(ast.NewIdentifierSelector(ast.NewIdentifier("x"))),
							ast.NewIntegerLiteral(3),
						),
					),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}
