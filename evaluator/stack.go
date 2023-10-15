package evaluator

import (
	"fmt"

	"github.com/rfejzic1/raiton/object"
)

type stack struct {
	values []object.Object
}

func (s *stack) push(o object.Object) {
	s.values = append(s.values, o)
}

func (s *stack) popSafe() (object.Object, error) {
	if len(s.values) == 0 {
		return nil, fmt.Errorf("expected a value on object stack")
	}

	return s.pop(), nil
}

func (s *stack) pop() object.Object {
	l := len(s.values)
	obj := s.values[l-1]
	s.values = s.values[:l-1]
	return obj
}
