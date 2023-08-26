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

func TestParser(t *testing.T) {
	source := `
	type person: {
	  name: string
	  age: number
	}
	`

	expected := Scope{
		Definitions: []*Definition{},
		TypeDefinitions: []*TypeDefinition{
			{
				Identifier: TypeIdentifier("person"),
				TypeExpression: &RecordType{
					Fields: map[Identifier]TypeExpression{
						Identifier("name"): NewTypeIdentifier("string"),
						Identifier("age"):  NewTypeIdentifier("number"),
					},
				},
			},
		},
		Expressions: []Expression{},
	}

	parseAndCompare(t, source, &expected)
}
