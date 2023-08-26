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
	main: 0
	`

	expected := Scope{
		Definitions:     make([]*Definition, 0),
		TypeDefinitions: make([]*TypeDefinition, 0),
		Expressions:     make([]Expression, 0),
	}

	parseAndCompare(t, source, &expected)
}
