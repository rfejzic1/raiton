package parser

import (
	"reflect"
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

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("assertion failed: ASTs are not equal")
	}
}

func TestParser(t *testing.T) {
	source := ``
	expected := Scope{
		Definitions:     make([]*Definition, 0),
		TypeDefinitions: make([]*TypeDefinition, 0),
		Expressions:     make([]Expression, 0),
	}

	parseAndCompare(t, source, &expected)
}
