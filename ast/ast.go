package ast

type Visitor interface {
	VisitScope(n *Scope) error
	VisitDefinition(n *Definition) error
	VisitIdentifier(n *Identifier) error
	VisitSelector(n *Selector) error
	VisitSelectorItem(n *SelectorItem) error
	VisitApplication(n *Application) error
	VisitFunction(n *FunctionLiteral) error
	VisitRecord(n *RecordLiteral) error
	VisitArray(n *ArrayLiteral) error
	VisitSlice(n *SliceLiteral) error
	VisitInteger(n *IntegerLiteral) error
	VisitFloat(n *FloatLiteral) error
	VisitString(n *StringLiteral) error
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

type Selector struct {
	Items []*SelectorItem
}

type SelectorItem struct {
	Identifier *Identifier
	Index      *IntegerLiteral
}

func NewIdentifierSelector(ident *Identifier) *SelectorItem {
	return &SelectorItem{
		Identifier: ident,
	}
}

func (i *SelectorItem) Accept(visitor Visitor) error {
	return visitor.VisitSelectorItem(i)
}

func NewIndexSelector(num *IntegerLiteral) *SelectorItem {
	return &SelectorItem{
		Index: num,
	}
}

func NewSelector(identifiers ...*SelectorItem) *Selector {
	return &Selector{
		Items: identifiers,
	}
}

func (f *Selector) Accept(visitor Visitor) error {
	return visitor.VisitSelector(f)
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

type IntegerLiteral int64

func NewIntegerLiteral(value int64) *IntegerLiteral {
	num := IntegerLiteral(value)
	return &num
}

func (n *IntegerLiteral) Accept(visitor Visitor) error {
	return visitor.VisitInteger(n)
}

type FloatLiteral float64

func NewFloatLiteral(value float64) *FloatLiteral {
	num := FloatLiteral(value)
	return &num
}

func (n *FloatLiteral) Accept(visitor Visitor) error {
	return visitor.VisitFloat(n)
}

type StringLiteral string

func NewStringLiteral(value string) *StringLiteral {
	string := StringLiteral(value)
	return &string
}

func (s *StringLiteral) Accept(visitor Visitor) error {
	return visitor.VisitString(s)
}

type BooleanLiteral string

func NewBooleanLiteral(value string) *BooleanLiteral {
	bool := BooleanLiteral(value)
	return &bool
}

func (b *BooleanLiteral) Accept(visitor Visitor) error {
	return visitor.VisitBoolean(b)
}
