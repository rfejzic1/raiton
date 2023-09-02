package ast

type Visitor interface {
	VisitScope(s *Scope) error
	VisitDefinition(d *Definition) error
	VisitTypeDefinition(d *TypeDefinition) error
	VisitTypeIdentifier(i *TypeIdentifier) error
	VisitTypeIdentifierPath(i *TypeIdentifierPath) error
	VisitFunctionType(f *FunctionType) error
	VisitRecordType(r *RecordType) error
	VisitSliceType(s *SliceType) error
	VisitArrayType(a *ArrayType) error
	VisitGroupType(g *GroupType) error
	VisitSumType(s *SumType) error
	VisitSumTypeVariant(v *SumTypeVariant) error
	VisitIdentifier(i *Identifier) error
	VisitIdentifierPath(i *IdentifierPath) error
	VisitInvocation(i *Invocation) error
	VisitLambda(l *LambdaLiteral) error
	VisitRecord(r *RecordLiteral) error
	VisitArray(a *ArrayLiteral) error
	VisitSlice(s *SliceLiteral) error
	VisitNumber(n *NumberLiteral) error
	VisitString(s *StringLiteral) error
	VisitCharacter(c *CharacterLiteral) error
}

type Node interface {
	Accept(visitor Visitor) error
}

type Scope struct {
	TypeDefinitions []*TypeDefinition
	Definitions     []*Definition
	Expressions     []Expression
}

func (s *Scope) Accept(visitor Visitor) error {
	return visitor.VisitScope(s)
}

type TypeDefinition struct {
	Identifier     TypeIdentifier
	Parameters     []*Identifier
	TypeExpression TypeExpression
}

func (d *TypeDefinition) Accept(visitor Visitor) error {
	return visitor.VisitTypeDefinition(d)
}

type Definition struct {
	TypeExpression TypeExpression
	Identifier     Identifier
	Parameters     []*Identifier
	Expression     Expression
}

func (d *Definition) Accept(visitor Visitor) error {
	return visitor.VisitDefinition(d)
}

type Expression interface {
	Node
}

type TypeExpression interface {
	Node
}

// *** Type Expressions ***

type TypeIdentifier string

func NewTypeIdentifier(name string) *TypeIdentifier {
	ident := TypeIdentifier(name)
	return &ident
}

func (i *TypeIdentifier) Accept(visitor Visitor) error {
	return visitor.VisitTypeIdentifier(i)
}

type TypeIdentifierPath struct {
	Identifiers []*TypeIdentifier
}

func NewTypeIdentifierPath(identifiers ...*TypeIdentifier) *TypeIdentifierPath {
	return &TypeIdentifierPath{
		Identifiers: identifiers,
	}
}

func (f *TypeIdentifierPath) Accept(visitor Visitor) error {
	return visitor.VisitTypeIdentifierPath(f)
}

type FunctionType struct {
	ParameterType TypeExpression
	ReturnType    TypeExpression
}

func (f *FunctionType) Accept(visitor Visitor) error {
	return visitor.VisitFunctionType(f)
}

type RecordType struct {
	Fields map[Identifier]TypeExpression
}

func (r *RecordType) Accept(visitor Visitor) error {
	return visitor.VisitRecordType(r)
}

type SliceType struct {
	ElementType TypeExpression
}

func (s *SliceType) Accept(visitor Visitor) error {
	return visitor.VisitSliceType(s)
}

type ArrayType struct {
	Size        uint64
	ElementType TypeExpression
}

func (a *ArrayType) Accept(visitor Visitor) error {
	return visitor.VisitArrayType(a)
}

type GroupType struct {
	TypeExpressions []TypeExpression
}

func (g *GroupType) Accept(visitor Visitor) error {
	return visitor.VisitGroupType(g)
}

type SumType struct {
	Variants []*SumTypeVariant
}

func (s *SumType) Accept(visitor Visitor) error {
	return visitor.VisitSumType(s)
}

type SumTypeVariant struct {
	Identifier     Identifier
	TypeExpression TypeExpression
}

func (v *SumTypeVariant) Accept(visitor Visitor) error {
	return visitor.VisitSumTypeVariant(v)
}

// *** Expressions ***

type Identifier string

func NewIdentifier(value string) *Identifier {
	ident := Identifier(value)
	return &ident
}

func (i *Identifier) Accept(visitor Visitor) error {
	return visitor.VisitIdentifier(i)
}

type IdentifierPath struct {
	Identifiers []*Identifier
}

func NewIdentifierPath(identifiers ...*Identifier) *IdentifierPath {
	return &IdentifierPath{
		Identifiers: identifiers,
	}
}

func (f *IdentifierPath) Accept(visitor Visitor) error {
	return visitor.VisitIdentifierPath(f)
}

type Invocation struct {
	Arguments []Expression
}

func NewInvocation(arguments ...Expression) *Invocation {
	return &Invocation{
		Arguments: arguments,
	}
}

func (i *Invocation) Accept(visitor Visitor) error {
	return visitor.VisitInvocation(i)
}

type LambdaLiteral struct {
	Parameters []*Identifier
	Expression Expression
}

func (l *LambdaLiteral) Accept(visitor Visitor) error {
	return visitor.VisitLambda(l)
}

type RecordLiteral struct {
	Fields map[Identifier]Expression
}

func (r *RecordLiteral) Accept(visitor Visitor) error {
	return visitor.VisitRecord(r)
}

type ArrayLiteral struct {
	Size     uint64
	Elements []Expression
}

func NewArrayLiteral(size uint64, elements ...Expression) *ArrayLiteral {
	return &ArrayLiteral{
		Size:     size,
		Elements: elements,
	}
}

func (a *ArrayLiteral) Accept(visitor Visitor) error {
	return visitor.VisitArray(a)
}

type SliceLiteral struct {
	Elements []Expression
}

func NewSliceLiteral(elements ...Expression) *SliceLiteral {
	return &SliceLiteral{
		Elements: elements,
	}
}

func (s *SliceLiteral) Accept(visitor Visitor) error {
	return visitor.VisitSlice(s)
}

type NumberLiteral string

func NewNumberLiteral(value string) *NumberLiteral {
	num := NumberLiteral(value)
	return &num
}

func (n *NumberLiteral) Accept(visitor Visitor) error {
	return visitor.VisitNumber(n)
}

type StringLiteral string

func NewStringLiteral(value string) *StringLiteral {
	string := StringLiteral(value)
	return &string
}

func (s *StringLiteral) Accept(visitor Visitor) error {
	return visitor.VisitString(s)
}

type CharacterLiteral string

func NewCharacterLiteral(value string) *CharacterLiteral {
	char := CharacterLiteral(value)
	return &char
}

func (c *CharacterLiteral) Accept(visitor Visitor) error {
	return visitor.VisitCharacter(c)
}
