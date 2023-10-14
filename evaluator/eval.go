package evaluator

import (
	"fmt"
	"strconv"

	"github.com/rfejzic1/raiton/ast"
	"github.com/rfejzic1/raiton/object"
)

type Evaluator struct {
	node    ast.Node
	results stack
}

func New(node ast.Node) Evaluator {
	return Evaluator{
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

var unsuported = fmt.Errorf("unsuported object")

func (e *Evaluator) VisitScope(s *ast.Scope) error {
	// TODO: Visit definitions

	var returnValue object.Object

	for _, expr := range s.Expressions {
		if err := expr.Accept(e); err != nil {
			return err
		}
		returnValue = e.results.pop()
	}

	// return the evaluation result of the last expression in scope
	e.results.push(returnValue)

	return nil
}

func (e *Evaluator) VisitDefinition(d *ast.Definition) error {
	return unsuported
}

func (e *Evaluator) VisitIdentifier(i *ast.Identifier) error {
	return unsuported
}

func (e *Evaluator) VisitIdentifierPath(i *ast.IdentifierPath) error {
	return unsuported
}

func (e *Evaluator) VisitApplication(i *ast.Application) error {
	return unsuported
}

func (e *Evaluator) VisitFunction(l *ast.FunctionLiteral) error {
	return unsuported
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
