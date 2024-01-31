package evaluator

import (
	"fmt"

	"raiton/ast"
	"raiton/object"
)

var builtins = map[string]*object.Builtin{
	"add": object.NewBuiltin(add),
	"eq":  object.NewBuiltin(eq),
	"map": object.NewBuiltin(mapfn),
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

func eq(_ ast.Visitor, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected two objects to compare")
	}

	first := args[0]
	second := args[1]

	if first.Type() != second.Type() {
		return nil, fmt.Errorf("expected both arguments to be of same type, but got %s and %s", first.Type(), second.Type())
	}

	switch f := first.(type) {
	case *object.Boolean:
		s := second.(*object.Boolean)
		return object.BoxBoolean(f.Value == s.Value), nil
	case *object.Integer:
		s := second.(*object.Integer)
		return object.BoxBoolean(f.Value == s.Value), nil
	case *object.Float:
		s := second.(*object.Float)
		return object.BoxBoolean(f.Value == s.Value), nil
	case *object.String:
		s := second.(*object.String)
		return object.BoxBoolean(f.Value == s.Value), nil
	default:
		return nil, fmt.Errorf("unsuported equality operation for type %s", f.Type())
	}
}

func mapfn(v ast.Visitor, args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected mapping function and array or list")
	}

	fn, ok := args[0].(*object.Function)

	if !ok {
		return nil, fmt.Errorf("expected first argument to be a function, but got %s", args[0].Type())
	}

	eval, ok := v.(*Evaluator)

	if !ok {
		return nil, fmt.Errorf("expected Evaluator visitor")
	}

	switch v := args[1].(type) {
	case *object.Array:
		newArray := &object.Array{
			Value: []object.Object{},
		}

		for _, arg := range v.Value {
			obj, err := eval.applyBuiltin(fn, arg)

			if err != nil {
				return nil, err
			}

			newArray.Value = append(newArray.Value, obj)
		}

		newArray.Size = uint64(len(newArray.Value))

		return newArray, nil
	case *object.List:
		newList := &object.List{
			Size: v.Size,
		}

		head := v.Head
		var newHead *object.ListNode

		for head != nil {
			value, err := eval.applyBuiltin(fn, head.Value)

			if err != nil {
				return nil, err
			}

			node := &object.ListNode{
				Value: value,
			}

			if newList.Head == nil {
				newHead = node
				newList.Head = newHead
			} else {
				newHead.Next = node
				newHead = node
			}

			head = head.Next
		}

		return newList, nil
	default:
		return nil, fmt.Errorf("expected second argument to be an array or list, but got %s", args[1].Type())
	}
}
