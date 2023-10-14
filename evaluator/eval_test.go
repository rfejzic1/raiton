package evaluator

import (
	"testing"

	"github.com/rfejzic1/raiton/lexer"
	"github.com/rfejzic1/raiton/object"
	"github.com/rfejzic1/raiton/parser"
)

func testEvaluation(env *object.Environment, input string) (object.Object, error) {
	l := lexer.New(input)
	p := parser.New(&l)
	program, err := p.Parse()

	if err != nil {
		return nil, err
	}

	eval := New(env, program)

	return eval.Evaluate()
}

func TestEvaluationInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	env := object.NewEnvironment()

	for _, tt := range tests {
		evaluated, err := testEvaluation(env, tt.input)

		if err != nil {
			t.Fatal(err)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}

}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not integer. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected %d, but got %d", result.Value, expected)
		return false
	}

	return true
}

func TestEvaluationBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	env := object.NewEnvironment()

	for _, tt := range tests {
		evaluated, err := testEvaluation(env, tt.input)

		if err != nil {
			t.Fatal(err)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not boolean. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. expected %t, but got %t", result.Value, expected)
		return false
	}

	return true
}
