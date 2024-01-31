package evaluator

import (
	"fmt"
	"strconv"

	"raiton/ast"
	"raiton/object"
)

type Evaluator struct {
	env      *object.Environment
	results  stack
	builtins builtins
}

func New(env *object.Environment) Evaluator {
	e := Evaluator{
		env: env,
	}

	e.builtins = newBuiltins(&e)

	return e
}

func (e *Evaluator) Evaluate(node ast.Node) (object.Object, error) {
	if err := e.evaluate(node); err != nil {
		return nil, err
	}

	return e.results.popSafe()
}

func (e *Evaluator) evaluate(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Scope:
		return e.scope(n)
	case *ast.Definition:
		return e.definition(n)
	case *ast.Identifier:
		return e.identifier(n)
	case *ast.Selector:
		return e.selector(n)
	case *ast.SelectorItem:
		return e.selectorItem(n)
	case *ast.Application:
		return e.application(n)
	case *ast.Function:
		return e.function(n)
	case *ast.Conditional:
		return e.conditional(n)
	case *ast.Record:
		return e.record(n)
	case *ast.Array:
		return e.array(n)
	case *ast.List:
		return e.list(n)
	case *ast.String:
		return e.string(n)
	case *ast.Integer:
		return e.integer(n)
	case *ast.Float:
		return e.float(n)
	case *ast.Keyword:
		return e.keyword(n)
	case *ast.Boolean:
		return e.boolean(n)
	default:
		panic("unhandled ast type")
	}
}

/*** Evaluator Methods ***/

func (e *Evaluator) scope(s *ast.Scope) error {
	for _, def := range s.Definitions {
		if err := e.evaluate(def); err != nil {
			return nil
		}
	}

	var returnValue object.Object

	for _, expr := range s.Expressions {
		if err := e.evaluate(expr); err != nil {
			return err
		}
		returnValue = e.results.pop()
	}

	// return the evaluation result of the last expression in scope
	if returnValue != nil {
		e.results.push(returnValue)
	}

	return nil
}

func (e *Evaluator) definition(d *ast.Definition) error {
	ident := string(d.Identifier)

	if err := e.evaluate(d.Expression); err != nil {
		return err
	}

	obj := e.results.pop()

	obj = e.env.Define(ident, obj)

	e.results.push(obj)

	return nil
}

func (e *Evaluator) identifier(i *ast.Identifier) error {
	ident := string(*i)

	if obj, ok := e.env.Lookup(ident); ok {
		e.results.push(obj)
		return nil
	}

	if obj, ok := e.builtins.lookup(ident); ok {
		e.results.push(obj)
		return nil
	}

	return fmt.Errorf("'%s' not defined", ident)
}

func (e *Evaluator) selector(s *ast.Selector) error {
	if len(s.Items) < 1 || s.Items[0].Identifier == nil {
		return fmt.Errorf("expected first selector item to be an identifier")
	}

	ident := string(*s.Items[0].Identifier)

	var obj object.Object
	var ok bool

	if obj, ok = e.env.Lookup(ident); !ok {
		if obj, ok = e.builtins.lookup(ident); !ok {
			return fmt.Errorf("'%s' not defined", ident)
		}
	}

	e.results.push(obj)

	for _, i := range s.Items[1:] {
		if err := e.evaluate(i); err != nil {
			return err
		}
	}

	return nil
}

func (e *Evaluator) selectorItem(i *ast.SelectorItem) error {
	obj := e.results.pop()

	switch obj.Type() {
	case object.RECORD:
		record := obj.(*object.Record)

		if i.Identifier == nil {
			return fmt.Errorf("can only access record fields with identifiers")
		}

		ident := string(*i.Identifier)

		obj, ok := record.Value[ident]

		if !ok {
			return fmt.Errorf("field '%s' not defined on record", ident)
		}

		e.results.push(obj)

		return nil
	case object.ARRAY:
		array := obj.(*object.Array)

		if i.Index == nil {
			return fmt.Errorf("can only access array elements with index")
		}

		index := int64(*i.Index)

		if index > int64(len(array.Value)) {
			return fmt.Errorf("index %d is out of bounds", index)
		}

		obj := array.Value[index]

		e.results.push(obj)

		return nil
	case object.LIST:
		list := obj.(*object.List)

		if i.Index == nil {
			return fmt.Errorf("can only access list elements with index")
		}

		index := int64(*i.Index)

		if index > int64(list.Size) {
			return fmt.Errorf("index %d is out of bounds", index)
		}

		head := list.Head
		counter := int64(0)

		for head != nil {
			if counter == index {
				e.results.push(head.Value)
				return nil
			}

			head = head.Next
		}

		return nil
	default:
		return fmt.Errorf("expected a collection but got %s", obj.Type())
	}
}

func (e *Evaluator) application(a *ast.Application) error {
	if len(a.Arguments) < 1 {
		e.results.push(object.THE_UNIT)
		return nil
	}

	if err := e.evaluate(a.Arguments[0]); err != nil {
		return err
	}

	obj := e.results.pop()

	switch obj.Type() {
	case object.FUNCTION:
		function := obj.(*object.Function)

		args := a.Arguments[1:]

		if len(args) < len(function.Parameters) {
			boundParams := function.Parameters[:len(args)]
			boundEnv := object.CloneEnvironment(function.Environment)

			newFunction := &object.Function{
				Parameters:  function.Parameters[len(args):],
				Body:        function.Body,
				Environment: boundEnv,
			}

			for i, arg := range args {
				param := boundParams[i]

				if err := e.evaluate(arg); err != nil {
					return err
				}

				ident := string(*param)
				obj := e.results.pop()

				newFunction.Environment.Define(ident, obj)
			}

			e.results.push(newFunction)

			return nil
		}

		e.env = object.NewEnclosedEnvironment(function.Environment)

		for i, p := range function.Parameters {
			arg := args[i]

			if err := e.evaluate(arg); err != nil {
				return err
			}

			ident := string(*p)
			obj := e.results.pop()

			e.env.Define(ident, obj)
		}

		if err := e.evaluate(function.Body); err != nil {
			return err
		}

		e.env = e.env.Enclosing()
	case object.BUILTIN:
		function := obj.(*object.Builtin).Fn

		objs := []object.Object{}

		// TODO: Partial function application here as well...
		for _, a := range a.Arguments[1:] {
			if err := e.evaluate(a); err != nil {
				return err
			}

			obj := e.results.pop()
			objs = append(objs, obj)
		}

		// TODO: Better error handling
		obj, err := function(objs...)

		if err != nil {
			return err
		}

		e.results.push(obj)
	default:
		e.results.push(obj)
	}

	return nil
}

func (e *Evaluator) applyBuiltin(fn *object.Function, args ...object.Object) (object.Object, error) {
	if len(args) != len(fn.Parameters) {
		return nil, fmt.Errorf("function expects %d arguments, but got %d", len(fn.Parameters), len(args))
	}

	for i, p := range fn.Parameters {
		ident := string(*p)
		e.env.Define(ident, args[i])
	}

	if err := e.evaluate(fn.Body); err != nil {
		return nil, err
	}

	obj := e.results.pop()
	return obj, nil
}

func (e *Evaluator) function(f *ast.Function) error {
	obj := &object.Function{
		Parameters:  f.Parameters,
		Body:        f.Body,
		Environment: object.NewEnclosedEnvironment(e.env),
	}

	e.results.push(obj)

	return nil
}

func (e *Evaluator) conditional(c *ast.Conditional) error {
	if err := e.evaluate(c.Condition); err != nil {
		return err
	}

	obj := e.results.pop()

	condition, ok := obj.(*object.Boolean)

	if !ok {
		return fmt.Errorf("expected a boolean value in if-expression condition")
	}

	if condition.Value {
		if err := e.evaluate(c.Consequence); err != nil {
			return nil
		}
	} else {
		if err := e.evaluate(c.Alternative); err != nil {
			return nil
		}
	}

	return nil
}

func (e *Evaluator) record(r *ast.Record) error {
	record := &object.Record{
		Value: map[string]object.Object{},
	}

	for field, value := range r.Fields {
		if err := e.evaluate(value); err != nil {
			return err
		}

		obj := e.results.pop()
		record.Value[string(field)] = obj
	}

	e.results.push(record)

	return nil
}

func (e *Evaluator) array(a *ast.Array) error {
	objs := []object.Object{}

	for _, elem := range a.Elements {
		if err := e.evaluate(elem); err != nil {
			return err
		}

		obj := e.results.pop()
		objs = append(objs, obj)
	}

	size := uint64(len(objs))

	if size != *a.Size {
		return fmt.Errorf("expected array of size %d, but got %d", a.Size, size)
	}

	array := &object.Array{
		Value: objs,
		Size:  size,
	}

	e.results.push(array)

	return nil
}

func (e *Evaluator) list(s *ast.List) error {
	list := &object.List{}

	var head *object.ListNode
	size := 0

	for _, elem := range s.Elements {
		if err := e.evaluate(elem); err != nil {
			return err
		}

		size += 1

		obj := e.results.pop()

		node := &object.ListNode{
			Value: obj,
		}

		if head == nil {
			head = node
			list.Head = head
			continue
		}

		head.Next = node
		head = node
	}

	e.results.push(list)

	return nil
}

func (e *Evaluator) integer(n *ast.Integer) error {
	result := &object.Integer{
		Value: int64(*n),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) float(n *ast.Float) error {
	result := &object.Float{
		Value: float64(*n),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) string(s *ast.String) error {
	result := &object.String{
		Value: string(*s),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) keyword(s *ast.Keyword) error {
	result := &object.Keyword{
		Value: string(*s),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) boolean(b *ast.Boolean) error {
	value, err := strconv.ParseBool(string(*b))

	if err != nil {
		return err
	}

	result := object.BoxBoolean(value)

	e.results.push(result)

	return nil
}
