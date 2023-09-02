package parser

import (
	"testing"

	"github.com/rfejzic1/raiton/lexer"
)

func parseAndCompare(t *testing.T, source string, expected Expression) {
	l := lexer.New(source)
	p := New(&l)
	got, err := p.Parse()

	if err != nil {
		t.Fatalf("parse error: %s", err)
	}

	comp := NewComparator(got)

	if err := comp.Compare(expected); err != nil {
		t.Fatalf("assertion failed: %s", err)
	}
}

func TestExpressionString(t *testing.T) {
	source := `"this is a string"`

	expected := Scope{
		Expressions: []Expression{
			NewStringLiteral("this is a string"),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionCharacter(t *testing.T) {
	source := `'c'`

	expected := Scope{
		Expressions: []Expression{
			NewCharacterLiteral("c"),
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

	expected := Scope{
		Expressions: []Expression{
			NewNumberLiteral("5"),
			NewNumberLiteral("2.65"),
			NewNumberLiteral("-1"),
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionArray(t *testing.T) {
	source := `[3: 1 2 3]`

	expected := Scope{
		Expressions: []Expression{
			&ArrayLiteral{
				Size: 3,
				Elements: []Expression{
					NewNumberLiteral("1"),
					NewNumberLiteral("2"),
					NewNumberLiteral("3"),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionSlice(t *testing.T) {
	source := `[1 2 3]`

	expected := Scope{
		Expressions: []Expression{
			&SliceLiteral{
				Elements: []Expression{
					NewNumberLiteral("1"),
					NewNumberLiteral("2"),
					NewNumberLiteral("3"),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestExpressionInvocation(t *testing.T) {
	source := `(println "Hello, World")`

	expected := Scope{
		Expressions: []Expression{
			&Invocation{
				Arguments: []Expression{
					NewIdentifierPath(NewIdentifier("println")),
					NewStringLiteral("Hello, World"),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestDefinitionTypedWithSingleExpression(t *testing.T) {
	source := `
	<string>
	name: "Tojuro"
	`

	expected := Scope{
		Definitions: []*Definition{
			{
				TypeExpression: NewTypeIdentifierPath(NewTypeIdentifier("string")),
				Identifier:     Identifier("name"),
				Expression:     NewStringLiteral("Tojuro"),
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

	expected := Scope{
		Definitions: []*Definition{
			{
				TypeExpression: NewTypeIdentifierPath(NewTypeIdentifier("number")),
				Identifier:     Identifier("age"),
				Expression: &Scope{
					Expressions: []Expression{
						NewNumberLiteral("24"),
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

	expected := Scope{
		Definitions: []*Definition{
			{
				TypeExpression: &FunctionType{
					ParameterType: NewTypeIdentifierPath(NewTypeIdentifier("number")),
					ReturnType:    NewTypeIdentifierPath(NewTypeIdentifier("number")),
				},
				Identifier: Identifier("add_two"),
				Parameters: []*Identifier{
					NewIdentifier("x"),
				},
				Expression: &Invocation{
					Arguments: []Expression{
						NewIdentifierPath(NewIdentifier("add")),
						NewIdentifierPath(NewIdentifier("x")),
						NewNumberLiteral("2"),
					},
				},
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

	expected := Scope{
		Definitions: []*Definition{
			{
				TypeExpression: &FunctionType{
					ParameterType: NewTypeIdentifierPath(NewTypeIdentifier("number")),
					ReturnType:    NewTypeIdentifierPath(NewTypeIdentifier("number")),
				},
				Identifier: Identifier("add_three"),
				Parameters: []*Identifier{
					NewIdentifier("x"),
				},
				Expression: &Scope{
					Expressions: []Expression{
						&Invocation{
							Arguments: []Expression{
								NewIdentifierPath(NewIdentifier("add")),
								NewIdentifierPath(NewIdentifier("x")),
								NewNumberLiteral("3"),
							},
						},
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier:     TypeIdentifier("name"),
				TypeExpression: NewTypeIdentifierPath(NewTypeIdentifier("string")),
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestTypeDefinitionArray(t *testing.T) {
	source := `
	type numArray: [3: number]
	`

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("numArray"),
				TypeExpression: &ArrayType{
					Size:        3,
					ElementType: NewTypeIdentifierPath(NewTypeIdentifier("number")),
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("numSlice"),
				TypeExpression: &SliceType{
					ElementType: NewTypeIdentifierPath(NewTypeIdentifier("number")),
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("person"),
				TypeExpression: &RecordType{
					Fields: map[Identifier]TypeExpression{
						Identifier("name"): NewTypeIdentifierPath(NewTypeIdentifier("string")),
						Identifier("age"):  NewTypeIdentifierPath(NewTypeIdentifier("number")),
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("color"),
				TypeExpression: &SumType{
					Variants: []*SumTypeVariant{
						{
							Identifier: Identifier("Red"),
						},
						{
							Identifier: Identifier("Green"),
						},
						{
							Identifier: Identifier("Blue"),
						},
						{
							Identifier: Identifier("RGB"),
							TypeExpression: &RecordType{
								Fields: map[Identifier]TypeExpression{
									Identifier("r"): NewTypeIdentifierPath(NewTypeIdentifier("number")),
									Identifier("g"): NewTypeIdentifierPath(NewTypeIdentifier("number")),
									Identifier("b"): NewTypeIdentifierPath(NewTypeIdentifier("number")),
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("option"),
				Parameters: []*Identifier{
					NewIdentifier("T"),
				},
				TypeExpression: &SumType{
					Variants: []*SumTypeVariant{
						{
							Identifier:     Identifier("Some"),
							TypeExpression: NewTypeIdentifierPath(NewTypeIdentifier("T")),
						},
						{
							Identifier: Identifier("None"),
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

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("stringOption"),
				TypeExpression: &GroupType{
					TypeExpressions: []TypeExpression{
						NewTypeIdentifierPath(NewTypeIdentifier("option")),
						NewTypeIdentifierPath(NewTypeIdentifier("string")),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}
