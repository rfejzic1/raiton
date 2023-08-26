package parser

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
	return expected.Accept(c)
}

func (c *Comparator) observe(node Node) {
	c.current = node
}

/*** Visitor Methods ***/

// NOTE: These visitor methods get called on the expectation tree.
//		 The c.current Node points to the matching node of
//		 the tree being checked and is updated and checked by the
//		 visitor methods.

func (c *Comparator) VisitScope(expected *Scope) error {
	current, ok := c.current.(*Scope)

	if !ok {
		return nodeTypeError("Scope")
	}

	if err := compareSlices(c, "definitions", expected.Definitions, current.Definitions); err != nil {
		return err
	}

	if err := compareSlices(c, "type definitions", expected.TypeDefinitions, current.TypeDefinitions); err != nil {
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

	c.observe(current.TypeExpression)

	if err := c.Compare(expected.TypeExpression); err != nil {
		return err
	}

	c.observe(&current.Identifier)

	if err := c.Compare(&expected.Identifier); err != nil {
		return err
	}

	if err := compareSlices(c, "parameters", expected.Parameters, current.Parameters); err != nil {
		return err
	}

	c.observe(current.Expression)

	if err := expected.Expression.Accept(c); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitTypeDefinition(expected *TypeDefinition) error {
	current, ok := c.current.(*TypeDefinition)

	if !ok {
		return nodeTypeError("TypeDefinition")
	}

	c.observe(&current.Identifier)

	if err := c.Compare(&expected.Identifier); err != nil {
		return err
	}

	c.observe(current.TypeExpression)

	if err := c.Compare(expected.TypeExpression); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitTypeIdentifier(expected *TypeIdentifier) error {
	current, ok := c.current.(*TypeIdentifier)

	if !ok {
		return nodeTypeError("TypeIdentifier")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) VisitFunctionType(expected *FunctionType) error {
	current, ok := c.current.(*FunctionType)

	if !ok {
		return nodeTypeError("FunctionType")
	}

	c.observe(current.ParameterType)

	if err := c.Compare(expected.ParameterType); err != nil {
		return err
	}

	c.observe(current.ReturnType)

	if err := c.Compare(expected.ReturnType); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitRecordType(expected *RecordType) error {
	current, ok := c.current.(*RecordType)

	if !ok {
		return nodeTypeError("RecordType")
	}

	if err := compareMaps(c, expected.Fields, current.Fields); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitSliceType(expected *SliceType) error {
	current, ok := c.current.(*SliceType)

	if !ok {
		return nodeTypeError("SliceType")
	}

	c.observe(current.ElementType)

	if err := c.Compare(expected.ElementType); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitArrayType(expected *ArrayType) error {
	current, ok := c.current.(*ArrayType)

	if !ok {
		return nodeTypeError("SliceType")
	}

	if expected.Size != current.Size {
		return fmt.Errorf("expected array of size %d, but got size %d", expected.Size, current.Size)
	}

	c.observe(current.ElementType)

	if err := c.Compare(expected.ElementType); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitGroupType(expected *GroupType) error {
	current, ok := c.current.(*GroupType)

	if !ok {
		return nodeTypeError("GroupType")
	}

	if err := compareSlices(c, "type expressions", expected.TypeExpressions, current.TypeExpressions); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitSumType(expected *SumType) error {
	current, ok := c.current.(*SumType)

	if !ok {
		return nodeTypeError("SumType")
	}

	if err := compareSlices(c, "variants", expected.Variants, current.Variants); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitSumTypeVariant(expected *SumTypeVariant) error {
	current, ok := c.current.(*SumTypeVariant)

	if !ok {
		return nodeTypeError("SumTypeVariant")
	}

	c.observe(&current.Identifier)

	if err := c.Compare(&expected.Identifier); err != nil {
		return err
	}

	c.observe(current.TypeExpression)

	if err := c.Compare(expected.TypeExpression); err != nil {
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

func (c *Comparator) VisitInvocation(expected *Invocation) error {
	current, ok := c.current.(*Invocation)

	if !ok {
		return nodeTypeError("Invocation")
	}

	if err := compareSlices(c, "arguments", expected.Arguments, current.Arguments); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitLambda(expected *LambdaLiteral) error {
	current, ok := c.current.(*LambdaLiteral)

	if !ok {
		return nodeTypeError("LambdaLiteral")
	}

	if err := compareSlices(c, "parameters", expected.Parameters, current.Parameters); err != nil {
		return err
	}

	c.observe(current.Expression)

	if err := c.Compare(expected.Expression); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitRecord(expected *RecordLiteral) error {
	current, ok := c.current.(*RecordLiteral)

	if !ok {
		return nodeTypeError("RecordLiteral")
	}

	if err := compareMaps(c, expected.Fields, current.Fields); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitArray(expected *ArrayLiteral) error {
	current, ok := c.current.(*ArrayLiteral)

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

func (c *Comparator) VisitSlice(expected *SliceLiteral) error {
	current, ok := c.current.(*SliceLiteral)

	if !ok {
		return nodeTypeError("SliceLiteral")
	}

	if err := compareSlices(c, "elements", expected.Elements, current.Elements); err != nil {
		return err
	}

	return nil
}

func (c *Comparator) VisitNumber(expected *NumberLiteral) error {
	current, ok := c.current.(*NumberLiteral)

	if !ok {
		return nodeTypeError("NumberLiteral")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) VisitString(expected *StringLiteral) error {
	current, ok := c.current.(*StringLiteral)

	if !ok {
		return nodeTypeError("StringLiteral")
	}

	if string(*current) != string(*expected) {
		return fmt.Errorf("expected `%s`, but got `%s`", string(*expected), string(*current))
	}

	return nil
}

func (c *Comparator) VisitCharacter(expected *CharacterLiteral) error {
	current, ok := c.current.(*CharacterLiteral)

	if !ok {
		return nodeTypeError("CharacterLiteral")
	}

	if current != expected {
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
			return nil
		}
	}
	return nil
}
