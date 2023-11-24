package ast

import "fmt"

type Comparator struct {
	current Node
}

// Creates a new Comparator with a Node to be compared as argument.
func NewComparator(compared Node) Comparator {
	return Comparator{
		current: compared,
	}
}

// Compares the Node given to the NewComparator constructor
// with the expected Node, returning an error if the equality
// comparison fails.
func (c *Comparator) Compare(expected Node) error {
	if expected == nil && c.current == nil {
		return nil
	}

	if expected == nil && c.current != nil {
		return fmt.Errorf("expected nil")
	}

	if expected != nil && c.current == nil {
		return fmt.Errorf("unexpected nil")
	}

	return expected.Accept(c)
}

func (c *Comparator) observe(node Node) {
	c.current = node
}

/*** Visitor Methods ***/

func (c *Comparator) VisitScope(expected *Scope) error {
	current, ok := c.current.(*Scope)

	if !ok {
		return nodeTypeError("Scope")
	}

	if err := compareSlices(c, "definitions", expected.Definitions, current.Definitions); err != nil {
		return err
	}

	if err := compareSlices(c, "expressions", expected.Expressions, current.Expressions); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitDefinition(expected *Definition) error {
	current, ok := c.current.(*Definition)

	if !ok {
		return nodeTypeError("Definition")
	}

	c.observe(&current.Identifier)

	if err := c.Compare(&expected.Identifier); err != nil {
		return err
	}

	c.observe(current.Expression)

	if err := expected.Expression.Accept(c); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitIdentifier(expected *Identifier) error {
	current, ok := c.current.(*Identifier)

	if !ok {
		return nodeTypeError("Identifier")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) VisitSelector(expected *Selector) error {
	current, ok := c.current.(*Selector)

	if !ok {
		return nodeTypeError("Selector")
	}

	if err := compareSlices(c, "identifiers", expected.Items, current.Items); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitSelectorItem(expected *SelectorItem) error {
	current, ok := c.current.(*SelectorItem)

	if !ok {
		return nodeTypeError("SelectorItem")
	}

	if current.Identifier != nil {
		c.observe(current.Identifier)
		if err := c.Compare(expected.Identifier); err != nil {
			return err
		}
	} else {
		c.observe(current.Index)
		if err := c.Compare(expected.Index); err != nil {
			return err
		}
	}

	return nil
}

func (c *Comparator) VisitApplication(expected *Application) error {
	current, ok := c.current.(*Application)

	if !ok {
		return nodeTypeError("Invocation")
	}

	if err := compareSlices(c, "arguments", expected.Arguments, current.Arguments); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitFunction(expected *Function) error {
	current, ok := c.current.(*Function)

	if !ok {
		return nodeTypeError("FunctionLiteral")
	}

	if err := compareSlices(c, "parameters", expected.Parameters, current.Parameters); err != nil {
		return err
	}

	c.observe(current.Body)

	if err := c.Compare(expected.Body); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitRecord(expected *Record) error {
	current, ok := c.current.(*Record)

	if !ok {
		return nodeTypeError("RecordLiteral")
	}

	if err := compareMaps(c, expected.Fields, current.Fields); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitArray(expected *Array) error {
	current, ok := c.current.(*Array)

	if !ok {
		return nodeTypeError("ArrayLiteral")
	}

	if expected.Size != current.Size {
		return fmt.Errorf("expected array of size %d, but got size %d", expected.Size, current.Size)
	}

	if err := compareSlices(c, "elements", expected.Elements, current.Elements); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitSlice(expected *Slice) error {
	current, ok := c.current.(*Slice)

	if !ok {
		return nodeTypeError("SliceLiteral")
	}

	if err := compareSlices(c, "elements", expected.Elements, current.Elements); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitInteger(expected *Integer) error {
	current, ok := c.current.(*Integer)

	if !ok {
		return nodeTypeError("NumberLiteral")
	}

	if *current != *expected {
		return fmt.Errorf("expected `%d`, but got `%d`", *expected, *current)
	}

	return nil
}

func (c *Comparator) VisitFloat(expected *Float) error {
	current, ok := c.current.(*Float)

	if !ok {
		return nodeTypeError("FloatLiteral")
	}

	if *current != *expected {
		return fmt.Errorf("expected `%g`, but got `%f`", *expected, *current)
	}

	return nil
}

func (c *Comparator) VisitBoolean(expected *Boolean) error {
	current, ok := c.current.(*Boolean)

	if !ok {
		return nodeTypeError("BooleanLiteral")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", *expected, *current)
	}

	return nil
}

func (c *Comparator) VisitString(expected *String) error {
	current, ok := c.current.(*String)

	if !ok {
		return nodeTypeError("StringLiteral")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

/*** Helper Functions ***/

func nodeTypeError(expected string) error {
	return fmt.Errorf("expected node of type `%s`", expected)
}

func compareSlices[T Node](c *Comparator, what string, expected []T, current []T) error {
	if len(expected) != len(current) {
		return fmt.Errorf("expected %d %s, but got %d", len(expected), what, len(current))
	}

	for i, d := range expected {
		c.observe(current[i])

		if err := c.Compare(d); err != nil {
			return err
		}
	}

	return nil
}

func compareMaps[T Node](c *Comparator, expected map[Identifier]T, current map[Identifier]T) error {
	if len(expected) != len(current) {
		return fmt.Errorf("expected %d fields, but got %d", len(expected), len(current))
	}

	for ident, expr := range expected {
		otherExpr, ok := current[ident]

		if !ok {
			return fmt.Errorf("field `%s` not found", ident)
		}

		c.observe(otherExpr)

		if err := c.Compare(expr); err != nil {
			return err
		}
	}
	return nil
}
