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

	return c.compare(expected)
}

func (c *Comparator) observe(node Node) {
	c.current = node
}

func (e *Comparator) compare(node Node) error {
	switch n := node.(type) {
	case *Scope:
		return e.scope(n)
	case *Definition:
		return e.definition(n)
	case *Identifier:
		return e.identifier(n)
	case *Selector:
		return e.selector(n)
	case *SelectorItem:
		return e.selectorItem(n)
	case *Application:
		return e.application(n)
	case *Function:
		return e.function(n)
	case *Conditional:
		return e.conditional(n)
	case *Record:
		return e.record(n)
	case *Array:
		return e.array(n)
	case *List:
		return e.list(n)
	case *String:
		return e.string(n)
	case *Integer:
		return e.integer(n)
	case *Float:
		return e.float(n)
	case *Keyword:
		return e.keyword(n)
	case *Boolean:
		return e.boolean(n)
	default:
		panic("unhandled ast type")
	}
}

/*** Comparator Methods ***/

func (c *Comparator) scope(expected *Scope) error {
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

func (c *Comparator) definition(expected *Definition) error {
	current, ok := c.current.(*Definition)

	if !ok {
		return nodeTypeError("Definition")
	}

	c.observe(&current.Identifier)

	if err := c.Compare(&expected.Identifier); err != nil {
		return err
	}

	c.observe(current.Expression)

	if err := c.compare(expected.Expression); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) identifier(expected *Identifier) error {
	current, ok := c.current.(*Identifier)

	if !ok {
		return nodeTypeError("Identifier")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) selector(expected *Selector) error {
	current, ok := c.current.(*Selector)

	if !ok {
		return nodeTypeError("Selector")
	}

	if err := compareSlices(c, "identifiers", expected.Items, current.Items); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) selectorItem(expected *SelectorItem) error {
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

func (c *Comparator) application(expected *Application) error {
	current, ok := c.current.(*Application)

	if !ok {
		return nodeTypeError("Invocation")
	}

	if err := compareSlices(c, "arguments", expected.Arguments, current.Arguments); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) conditional(expected *Conditional) error {
	current, ok := c.current.(*Conditional)

	if !ok {
		return nodeTypeError("Conditional")
	}

	c.observe(current.Condition)

	if err := c.Compare(expected.Condition); err != nil {
		return err
	}

	c.observe(current.Consequence)

	if err := c.Compare(expected.Consequence); err != nil {
		return err
	}

	c.observe(current.Alternative)

	if err := c.Compare(expected.Alternative); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) function(expected *Function) error {
	current, ok := c.current.(*Function)

	if !ok {
		return nodeTypeError("Function")
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

func (c *Comparator) record(expected *Record) error {
	current, ok := c.current.(*Record)

	if !ok {
		return nodeTypeError("Record")
	}

	if err := compareMaps(c, expected.Fields, current.Fields); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) array(expected *Array) error {
	current, ok := c.current.(*Array)

	if !ok {
		return nodeTypeError("Array")
	}

	if *expected.Size != *current.Size {
		return fmt.Errorf("expected array of size %d, but got size %d", expected.Size, current.Size)
	}

	if err := compareSlices(c, "elements", expected.Elements, current.Elements); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) list(expected *List) error {
	current, ok := c.current.(*List)

	if !ok {
		return nodeTypeError("List")
	}

	if err := compareSlices(c, "elements", expected.Elements, current.Elements); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) integer(expected *Integer) error {
	current, ok := c.current.(*Integer)

	if !ok {
		return nodeTypeError("Number")
	}

	if *current != *expected {
		return fmt.Errorf("expected `%d`, but got `%d`", *expected, *current)
	}

	return nil
}

func (c *Comparator) float(expected *Float) error {
	current, ok := c.current.(*Float)

	if !ok {
		return nodeTypeError("Float")
	}

	if *current != *expected {
		return fmt.Errorf("expected `%g`, but got `%f`", *expected, *current)
	}

	return nil
}

func (c *Comparator) boolean(expected *Boolean) error {
	current, ok := c.current.(*Boolean)

	if !ok {
		return nodeTypeError("Boolean")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", *expected, *current)
	}

	return nil
}

func (c *Comparator) string(expected *String) error {
	current, ok := c.current.(*String)

	if !ok {
		return nodeTypeError("String")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) keyword(expected *Keyword) error {
	current, ok := c.current.(*Keyword)

	if !ok {
		return nodeTypeError("Keyword")
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
