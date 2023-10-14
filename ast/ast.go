package ast

type Visitor interface {
	VisitScope(s *Scope) error
	VisitDefinition(d *Definition) error
	VisitIdentifier(i *Identifier) error
	VisitIdentifierPath(i *IdentifierPath) error
	VisitApplication(i *Application) error
	VisitFunction(l *FunctionLiteral) error
	VisitRecord(r *RecordLiteral) error
	VisitArray(a *ArrayLiteral) error
	VisitSlice(s *SliceLiteral) error
	VisitNumber(n *NumberLiteral) error
	VisitString(s *StringLiteral) error
	VisitCharacter(c *CharacterLiteral) error
	VisitBoolean(c *BooleanLiteral) error
}

type Node interface {
	Accept(visitor Visitor) error
}

type Scope struct {
	Definitions []*Definition
	Expressions []Expression
}

func (s *Scope) Accept(visitor Visitor) error {
	return visitor.VisitScope(s)
}

type Definition struct {
	Identifier Identifier
	Parameters []*Identifier
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
	Expression Expression
}

func (l *FunctionLiteral) Accept(visitor Visitor) error {
	return visitor.VisitFunction(l)
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
