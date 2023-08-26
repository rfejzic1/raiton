package parser

type Visitor interface {
	VisitScope(s *Scope) error
	VisitDefinition(d *Definition) error
	VisitTypeDefinition(d *TypeDefinition) error
	VisitTypeIdentifier(i *TypeIdentifier) error
	VisitFunctionType(f *FunctionType) error
	VisitRecordType(r *RecordType) error
	VisitSliceType(s *SliceType) error
	VisitArrayType(a *ArrayType) error
	VisitGroupType(g *GroupType) error
	VisitSumType(s *SumType) error
	VisitSumTypeVariant(v *SumTypeVariant) error
	VisitIdentifier(i *Identifier) error
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
	typeDefinitions []TypeDefinition
	definitions     []Definition
	expressions     []Expression
}

func (s *Scope) Accept(visitor Visitor) error {
	return visitor.VisitScope(s)
}

type TypeDefinition struct {
	identifier     TypeIdentifier
	typeExpression TypeExpression
}

func (d *TypeDefinition) Accept(visitor Visitor) error {
	return visitor.VisitTypeDefinition(d)
}

type Definition struct {
	typeExpression TypeExpression
	identifier     Identifier
	parameters     []Identifier
	expression     Expression
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

func (i *TypeIdentifier) Accept(visitor Visitor) error {
	return visitor.VisitTypeIdentifier(i)
}

type FunctionType struct {
	parameterType TypeExpression
	returnType    TypeExpression
}

func (f *FunctionType) Accept(visitor Visitor) error {
	return visitor.VisitFunctionType(f)
}

type RecordType struct {
	fields map[Identifier]TypeExpression
}

func (r *RecordType) Accept(visitor Visitor) error {
	return visitor.VisitRecordType(r)
}

type SliceType struct {
	elementType TypeExpression
}

func (s *SliceType) Accept(visitor Visitor) error {
	return visitor.VisitSliceType(s)
}

type ArrayType struct {
	size        uint64
	elementType TypeExpression
}

func (a *ArrayType) Accept(visitor Visitor) error {
	return visitor.VisitArrayType(a)
}

type GroupType struct {
	typeExpressions []TypeExpression
}

func (g *GroupType) Accept(visitor Visitor) error {
	return visitor.VisitGroupType(g)
}

type SumType struct {
	variants []SumTypeVariant
}

func (s *SumType) Accept(visitor Visitor) error {
	return visitor.VisitSumType(s)
}

type SumTypeVariant struct {
	identifier     Identifier
	typeExpression TypeExpression
}

func (v *SumTypeVariant) Accept(visitor Visitor) error {
	return visitor.VisitSumTypeVariant(v)
}

// *** Expressions ***

type Identifier string

func (i *Identifier) Accept(visitor Visitor) error {
	return visitor.VisitIdentifier(i)
}

type Invocation struct {
	arguments []Expression
}

func (i *Invocation) Accept(visitor Visitor) error {
	return visitor.VisitInvocation(i)
}

type LambdaLiteral struct {
	parameters []Identifier
	expression Expression
}

func (l *LambdaLiteral) Accept(visitor Visitor) error {
	return visitor.VisitLambda(l)
}

type RecordLiteral struct {
	fields map[Identifier]Expression
}

func (r *RecordLiteral) Accept(visitor Visitor) error {
	return visitor.VisitRecord(r)
}

type ArrayLiteral struct {
	size     uint64
	elements []Expression
}

func (a *ArrayLiteral) Accept(visitor Visitor) error {
	return visitor.VisitArray(a)
}

type SliceLiteral struct {
	elements []Expression
}

func (s *SliceLiteral) Accept(visitor Visitor) error {
	return visitor.VisitSlice(s)
}

type NumberLiteral string

func (n *NumberLiteral) Accept(visitor Visitor) error {
	return visitor.VisitNumber(n)
}

type StringLiteral string

func (s *StringLiteral) Accept(visitor Visitor) error {
	return visitor.VisitString(s)
}

type CharacterLiteral string

func (c *CharacterLiteral) Accept(visitor Visitor) error {
	return visitor.VisitCharacter(c)
}
