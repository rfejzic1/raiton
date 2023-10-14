package evaluator

import (
	"fmt"

	"github.com/rfejzic1/raiton/object"
)

var builtins = map[string]*object.Builtin{
	"add": object.MakeBuiltin(func(args ...object.Object) (object.Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("expected two integers")
		}

		first, ok := args[0].(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("expected first argument to be integer, but got %s", args[0].Type())
		}

		second, ok := args[1].(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("expected second argument to be integer, but got %s", args[0].Type())
		}

		result := first.Value + second.Value

		return &object.Integer{
			Value: result,
		}, nil
	}),
}
