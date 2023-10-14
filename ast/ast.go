package ast

type Visitor interface {
	VisitScope(n *Scope) error
	VisitDefinition(n *Definition) error
	VisitIdentifier(n *Identifier) error
	VisitIdentifierPath(n *IdentifierPath) error
	VisitApplication(n *Application) error
	VisitFunction(n *FunctionLiteral) error
	VisitRecord(n *RecordLiteral) error
	VisitArray(n *ArrayLiteral) error
	VisitSlice(n *SliceLiteral) error
	VisitNumber(n *NumberLiteral) error
	VisitString(n *StringLiteral) error
	VisitCharacter(n *CharacterLiteral) error
	VisitBoolean(n *BooleanLiteral) error
}

type Node interface {
	Accept(visitor Visitor) error
}

type Scope struct {
	Definitions []*Definition
	Expressions []Expression
}

func ScopeExpressions(expressions ...Expression) *Scope {
	return &Scope{
		Definitions: []*Definition{},
		Expressions: expressions,
	}
}

func (s *Scope) Accept(visitor Visitor) error {
	return visitor.VisitScope(s)
}

type Definition struct {
	Identifier Identifier
	Expression Expression
}

func (d *Definition) Accept(visitor Visitor) error {
	return visitor.VisitDefinition(d)
}

// *** Expressions ***

type Expression interface {
	Node
}

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

type Application struct {
	Arguments []Expression
}

func NewApplication(arguments ...Expression) *Application {
	return &Application{
		Arguments: arguments,
	}
}

func (i *Application) Accept(visitor Visitor) error {
	return visitor.VisitApplication(i)
}

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *Scope
}

func (f *FunctionLiteral) Accept(visitor Visitor) error {
	return visitor.VisitFunction(f)
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

type BooleanLiteral string

func NewBooleanLiteral(value string) *BooleanLiteral {
	bool := BooleanLiteral(value)
	return &bool
}

func (b *BooleanLiteral) Accept(visitor Visitor) error {
	return visitor.VisitBoolean(b)
}
