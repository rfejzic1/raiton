package evaluator

import (
	"fmt"
	"strconv"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/object"
)

type Evaluator struct {
	env     *object.Environment
	node    ast.Node
	results stack
}

func New(env *object.Environment, node ast.Node) Evaluator {
	return Evaluator{
		env:  env,
		node: node,
	}
}

func (e *Evaluator) Evaluate() (object.Object, error) {
	if err := e.node.Accept(e); err != nil {
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

	obj, ok := e.env.Lookup(ident)

	if !ok {
		return fmt.Errorf("'%s' not defined", ident)
	}

	e.results.push(obj)

	return nil
}

func (e *Evaluator) VisitIdentifierPath(i *ast.IdentifierPath) error {
	// TODO: Only records and modules support identifier paths

	ident := string(*i.Identifiers[0])

	obj, ok := e.env.Lookup(ident)

	if !ok {
		return fmt.Errorf("'%s' not defined", ident)
	}

	e.results.push(obj)

	return nil
}

func (e *Evaluator) VisitApplication(a *ast.Application) error {
	if len(a.Arguments) < 1 {
		return fmt.Errorf("expected function")
	}

	if err := a.Arguments[0].Accept(e); err != nil {
		return err
	}

	obj := e.results.pop()

	if obj.Type() != object.FUNCTION {
		return fmt.Errorf("expected a function, but got %s", obj.Type())
	}

	function := obj.(*object.Function)

	args := a.Arguments[1:]

	if len(args) != len(function.Parameters) {
		return fmt.Errorf("function expects %d arguments, but got %d", len(function.Parameters), len(args))
	}

	for i, p := range function.Parameters {
		arg := args[i]

		if err := arg.Accept(e); err != nil {
			return err
		}

		ident := string(*p)
		obj := e.results.pop()

		// TODO: This should be added to new env
		e.env.Define(ident, obj)
	}

	if err := function.Body.Accept(e); err != nil {
		return err
	}

	return nil
}

func (e *Evaluator) VisitFunction(f *ast.FunctionLiteral) error {
	obj := &object.Function{
		Parameters:  f.Parameters,
		Body:        f.Body,
		Environment: e.env,
	}

	e.results.push(obj)

	return nil
}

func (e *Evaluator) VisitRecord(r *ast.RecordLiteral) error {
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

func (e *Evaluator) VisitArray(a *ast.ArrayLiteral) error {
	objs := []object.Object{}

	for _, elem := range a.Elements {
		if err := elem.Accept(e); err != nil {
			return err
		}

		obj := e.results.pop()
		objs = append(objs, obj)
	}

	size := uint64(len(objs))

	if size != a.Size {
		return fmt.Errorf("expected array of size %d, but got %d", a.Size, size)
	}

	array := &object.Array{
		Value: objs,
		Size:  size,
	}

	e.results.push(array)

	return nil
}

func (e *Evaluator) VisitSlice(s *ast.SliceLiteral) error {
	objs := []object.Object{}

	for _, elem := range s.Elements {
		if err := elem.Accept(e); err != nil {
			return err
		}

		obj := e.results.pop()
		objs = append(objs, obj)
	}

	size := uint64(len(objs))

	array := &object.Array{
		Value: objs,
		Size:  size,
	}

	slice := &object.Slice{
		Value: array,
	}

	e.results.push(slice)

	return nil
}

func (e *Evaluator) VisitNumber(n *ast.NumberLiteral) error {
	// TODO: Floats
	value, err := strconv.ParseInt(string(*n), 0, 64)

	if err != nil {
		return err
	}

	result := &object.Integer{
		Value: value,
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitString(s *ast.StringLiteral) error {
	result := &object.String{
		Value: string(*s),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitCharacter(c *ast.CharacterLiteral) error {
	result := &object.Character{
		Value: string(*c),
	}

	e.results.push(result)

	return nil
}

func (e *Evaluator) VisitBoolean(b *ast.BooleanLiteral) error {
	value, err := strconv.ParseBool(string(*b))

	if err != nil {
		return err
	}

	result := object.BoxBoolean(value)

	e.results.push(result)

	return nil
}
