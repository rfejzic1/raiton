package evaluator

import (
	"fmt"
	"strconv"

	"raiton/ast"
	"raiton/object"
)

type Evaluator struct {
	env     *object.Environment
	results stack
}

func New(env *object.Environment) Evaluator {
	return Evaluator{
		env: env,
	}
}

func (e *Evaluator) Evaluate(node ast.Node) (object.Object, error) {
	if err := node.Accept(e); err != nil {
		return nil, err
	}

	return e.results.popSafe()
}

/*** Visitor Methods ***/

func (e *Evaluator) VisitScope(s *ast.Scope) error {
	for _, def := range s.Definitions {
		if err := def.Accept(e); err != nil {
			return nil
		}
	}

	var returnValue object.Object

	for _, expr := range s.Expressions {
		if err := expr.Accept(e); err != nil {
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

func (e *Evaluator) VisitDefinition(d *ast.Definition) error {
	ident := string(d.Identifier)

	if err := d.Expression.Accept(e); err != nil {
		return err
	}

	obj := e.results.pop()

	obj = e.env.Define(ident, obj)

	e.results.push(obj)

	return nil
}

func (e *Evaluator) VisitIdentifier(i *ast.Identifier) error {
	ident := string(*i)

	if obj, ok := e.env.Lookup(ident); ok {
		e.results.push(obj)
		return nil
	}

	if obj, ok := builtins[ident]; ok {
		e.results.push(obj)
		return nil
	}

	return fmt.Errorf("'%s' not defined", ident)
}

func (e *Evaluator) VisitSelector(s *ast.Selector) error {
	if len(s.Items) < 1 || s.Items[0].Identifier == nil {
		return fmt.Errorf("expected first selector item to be an identifier")
	}

	ident := string(*s.Items[0].Identifier)

	var obj object.Object
	var ok bool

	if obj, ok = e.env.Lookup(ident); !ok {
		if obj, ok = builtins[ident]; !ok {
			return fmt.Errorf("'%s' not defined", ident)
		}
	}

	e.results.push(obj)

	for _, i := range s.Items[1:] {
		if err := i.Accept(e); err != nil {
			return err
		}
	}

	return nil
}

func (e *Evaluator) VisitSelectorItem(i *ast.SelectorItem) error {
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

func (e *Evaluator) VisitApplication(a *ast.Application) error {
	if len(a.Arguments) < 1 {
		return fmt.Errorf("expected at least one expression")
	}

	if err := a.Arguments[0].Accept(e); err != nil {
		return err
	}

	obj := e.results.pop()

	switch obj.Type() {
	case object.FUNCTION:
		function := obj.(*object.Function)

		args := a.Arguments[1:]

		if len(args) != len(function.Parameters) {
			return fmt.Errorf("function expects %d arguments, but got %d", len(function.Parameters), len(args))
		}

		e.env = object.NewEnclosedEnvironment(e.env)

		for i, p := range function.Parameters {
			arg := args[i]

			if err := arg.Accept(e); err != nil {
				return err
			}

			ident := string(*p)
			obj := e.results.pop()

			e.env.Define(ident, obj)
		}

		if err := function.Body.Accept(e); err != nil {
			return err
		}

		e.env = e.env.Enclosing()
	case object.BUILTIN:
		function := obj.(*object.Builtin).Fn

		objs := []object.Object{}

		for _, a := range a.Arguments[1:] {
			if err := a.Accept(e); err != nil {
				return err
			}

			obj := e.results.pop()
			objs = append(objs, obj)
		}

		// TODO: Better error handling
		obj, err := function(e, objs...)

		if err != nil {
			return err
		}

		e.results.push(obj)
	default:
		e.results.push(obj)
	}

	return nil
}

func (e *Evaluator) applyFunction(fn *object.Function, args ...object.Object) (object.Object, error) {
	if len(args) != len(fn.Parameters) {
		return nil, fmt.Errorf("function expects %d arguments, but got %d", len(fn.Parameters), len(args))
	}

	for i, p := range fn.Parameters {
		ident := string(*p)
		e.env.Define(ident, args[i])
	}

	if err := fn.Body.Accept(e); err != nil {
		return nil, err
	}

	obj := e.results.pop()
	return obj, nil
}

func (e *Evaluator) VisitFunction(f *ast.Function) error {
	obj := &object.Function{
		Parameters:  f.Parameters,
		Body:        f.Body,
		Environment: e.env,
	}

	e.results.push(obj)

	return nil
}

func (e *Evaluator) VisitRecord(r *ast.Record) error {
	record := &object.Record{
		Value: map[string]object.Object{},
	}

	for field, value := range r.Fields {
		if err := value.Accept(e); err != nil {
			return err
		}

		obj := e.results.pop()
		record.Value[string(field)] = obj
	}

	e.results.push(record)

	return nil
}

func (e *Evaluator) VisitArray(a *ast.Array) error {
	objs := []object.Object{}

	for _, elem := range a.Elements {
		if err := elem.Accept(e); err != nil {
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

func (e *Evaluator) VisitList(s *ast.List) error {
	list := &object.List{}

	var head *object.ListNode
	size := 0

	for _, elem := range s.Elements {
		if err := elem.Accept(e); err != nil {
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

func (e *Evaluator) VisitInteger(n *ast.Integer) error {
	result := &object.Integer{
		Value: int64(*n),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitFloat(n *ast.Float) error {
	result := &object.Float{
		Value: float64(*n),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitString(s *ast.String) error {
	result := &object.String{
		Value: string(*s),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitBoolean(b *ast.Boolean) error {
	value, err := strconv.ParseBool(string(*b))

	if err != nil {
		return err
	}

	result := object.BoxBoolean(value)

	e.results.push(result)

	return nil
}
