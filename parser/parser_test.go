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

func TestExpressions(t *testing.T) {
	source := `
	# string expression
	"string"

	# character expression
	'c'

	# number expression; positive integer
	5

	# number expression; positive float
	2.65

	# number expression; negative integer
	-1

	# array expression
	[3: 1 2 3]

	# slice expression
	[1 2 3]

	(println "Hello, World")
	`

	expected := Scope{
		Expressions: []Expression{
			NewStringLiteral("string"),
			NewCharacterLiteral("c"),
			NewNumberLiteral("5"),
			NewNumberLiteral("2.65"),
			NewNumberLiteral("-1"),
			&ArrayLiteral{
				Size: 3,
				Elements: []Expression{
					NewNumberLiteral("1"),
					NewNumberLiteral("2"),
					NewNumberLiteral("3"),
				},
			},
			&SliceLiteral{
				Elements: []Expression{
					NewNumberLiteral("1"),
					NewNumberLiteral("2"),
					NewNumberLiteral("3"),
				},
			},
			&Invocation{
				Arguments: []Expression{
					NewIdentifier("println"),
					NewStringLiteral("Hello, World"),
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}

func TestDefinition(t *testing.T) {
	source := `
	# definition
	<string>
	name: "Tojuro"

	# definition with scope
	<number>
	age { 24 }

	# function definition
	<number -> number>
	add_two x: (add x 2)

	# function definition with scope
	<number -> number>
	add_three x { (add x 3) }
	`

	expected := Scope{
		Expressions: []Expression{},
		Definitions: []*Definition{
			{
				TypeExpression: NewTypeIdentifier("string"),
				Identifier:     Identifier("name"),
				Expression:     NewStringLiteral("Tojuro"),
			},
			{
				TypeExpression: NewTypeIdentifier("number"),
				Identifier:     Identifier("age"),
				Expression: &Scope{
					Expressions: []Expression{
						NewNumberLiteral("24"),
					},
				},
			},
			{
				TypeExpression: &FunctionType{
					ParameterType: NewTypeIdentifier("number"),
					ReturnType:    NewTypeIdentifier("number"),
				},
				Identifier: Identifier("add_two"),
				Parameters: []*Identifier{
					NewIdentifier("x"),
				},
				Expression: &Invocation{
					Arguments: []Expression{
						NewIdentifier("add"),
						NewIdentifier("x"),
						NewNumberLiteral("2"),
					},
				},
			},
			{
				TypeExpression: &FunctionType{
					ParameterType: NewTypeIdentifier("number"),
					ReturnType:    NewTypeIdentifier("number"),
				},
				Identifier: Identifier("add_three"),
				Parameters: []*Identifier{
					NewIdentifier("x"),
				},
				Expression: &Scope{
					Expressions: []Expression{
						&Invocation{
							Arguments: []Expression{
								NewIdentifier("add"),
								NewIdentifier("x"),
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

func TestTypeDefinition(t *testing.T) {
	source := `
	type name: string

	type numArray: [3: number]

	type numSlice: [number]

	type person: {
		name: string
		age: number
	}

	type color: | Red | Green | Blue | RGB: { r:number g:number b:number }

	type option T: | Some: T | None

	type stringOption: (option string)
	`

	expected := Scope{
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier:     TypeIdentifier("name"),
				TypeExpression: NewTypeIdentifier("string"),
			},
			{
				Identifier: TypeIdentifier("numArray"),
				TypeExpression: &ArrayType{
					Size:        3,
					ElementType: NewTypeIdentifier("number"),
				},
			},
			{
				Identifier: TypeIdentifier("numSlice"),
				TypeExpression: &SliceType{
					ElementType: NewTypeIdentifier("number"),
				},
			},
			{
				Identifier: TypeIdentifier("person"),
				TypeExpression: &RecordType{
					Fields: map[Identifier]TypeExpression{
						Identifier("name"): NewTypeIdentifier("string"),
						Identifier("age"):  NewTypeIdentifier("number"),
					},
				},
			},
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
									Identifier("r"): NewTypeIdentifier("number"),
									Identifier("g"): NewTypeIdentifier("number"),
									Identifier("b"): NewTypeIdentifier("number"),
								},
							},
						},
					},
				},
			},
			{
				Identifier: TypeIdentifier("option"),
				Parameters: []*Identifier{
					NewIdentifier("T"),
				},
				TypeExpression: &SumType{
					Variants: []*SumTypeVariant{
						{
							Identifier:     Identifier("Some"),
							TypeExpression: NewTypeIdentifier("T"),
						},
						{
							Identifier: Identifier("None"),
						},
					},
				},
			},
			{
				Identifier: TypeIdentifier("stringOption"),
				TypeExpression: &GroupType{
					TypeExpressions: []TypeExpression{
						NewTypeIdentifier("option"),
						NewTypeIdentifier("string"),
					},
				},
			},
		},
	}

	parseAndCompare(t, source, &expected)
}
