package evaluator

import (
	"fmt"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/object"
)

var builtins = map[string]*object.Builtin{
	"add": object.MakeBuiltin(add),
	"map": object.MakeBuiltin(mapfn),
}

func add(_ ast.Visitor, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected two integers")
	}

	first, ok := args[0].(*object.Integer)
	if !ok {
		return nil, fmt.Errorf("expected first argument to be integer, but got %s", args[0].Type())
	}

	second, ok := args[1].(*object.Integer)
	if !ok {
		return nil, fmt.Errorf("expected second argument to be integer, but got %s", args[1].Type())
	}

	result := first.Value + second.Value

	return &object.Integer{
		Value: result,
	}, nil
}

func mapfn(v ast.Visitor, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected array and mapping function")
	}

	arr, ok := args[0].(*object.Array)

	if !ok {
		if slc, ok := args[0].(*object.Slice); ok {
			arr = slc.Value
		} else {
			return nil, fmt.Errorf("expected first argument to be an array, but got %s", args[0].Type())
		}
	}

	fn, ok := args[1].(*object.Function)
	if !ok {
		return nil, fmt.Errorf("expected second argument to be a function, but got %s", args[1].Type())
	}

	eval, ok := v.(*Evaluator)

	if !ok {
		return nil, fmt.Errorf("expected Evaluator visitor")
	}

	newArray := &object.Array{
		Value: []object.Object{},
	}

	for _, arg := range arr.Value {
		obj, err := eval.applyFunction(fn, arg)
		if err != nil {
			return nil, err
		}

		newArray.Value = append(newArray.Value, obj)
	}

	newArray.Size = uint64(len(newArray.Value))

	return newArray, nil
}
