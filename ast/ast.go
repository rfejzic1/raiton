package ast

type Visitor interface {
	VisitScope(n *Scope) error
	VisitDefinition(n *Definition) error
	VisitIdentifier(n *Identifier) error
	VisitSelector(n *Selector) error
	VisitSelectorItem(n *SelectorItem) error
	VisitApplication(n *Application) error
	VisitFunction(n *Function) error
	VisitConditional(n *Conditional) error
	VisitRecord(n *Record) error
	VisitArray(n *Array) error
	VisitList(n *List) error
	VisitInteger(n *Integer) error
	VisitFloat(n *Float) error
	VisitString(n *String) error
	VisitKeyword(n *Keyword) error
	VisitBoolean(n *Boolean) error
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
	Index      *Integer
}

func NewIdentifierSelector(ident *Identifier) *SelectorItem {
	return &SelectorItem{
		Identifier: ident,
	}
}

func (i *SelectorItem) Accept(visitor Visitor) error {
	return visitor.VisitSelectorItem(i)
}

func NewIndexSelector(num *Integer) *SelectorItem {
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

type Function struct {
	Parameters []*Identifier
	Body       *Scope
}

func (f *Function) Accept(visitor Visitor) error {
	return visitor.VisitFunction(f)
}

type Conditional struct {
	Condition   Expression
	Consequence *Scope
	Alternative *Scope
}

func (c *Conditional) Accept(visitor Visitor) error {
	return visitor.VisitConditional(c)
}

type Record struct {
	Fields map[Identifier]Expression
}

func (r *Record) Accept(visitor Visitor) error {
	return visitor.VisitRecord(r)
}

type Array struct {
	Size     *uint64
	Elements []Expression
}

func NewArray(size uint64, elements ...Expression) *Array {
	return &Array{
		Size:     &size,
		Elements: elements,
	}
}

func (a *Array) Accept(visitor Visitor) error {
	return visitor.VisitArray(a)
}

type List struct {
	Elements []Expression
}

func NewList(elements ...Expression) *List {
	return &List{
		Elements: elements,
	}
}

func (s *List) Accept(visitor Visitor) error {
	return visitor.VisitList(s)
}

type Integer int64

func NewInteger(value int64) *Integer {
	i := Integer(value)
	return &i
}

func (n *Integer) Accept(visitor Visitor) error {
	return visitor.VisitInteger(n)
}

type Float float64

func NewFloat(value float64) *Float {
	f := Float(value)
	return &f
}

func (n *Float) Accept(visitor Visitor) error {
	return visitor.VisitFloat(n)
}

type String string

func NewString(value string) *String {
	s := String(value)
	return &s
}

func (s *String) Accept(visitor Visitor) error {
	return visitor.VisitString(s)
}

type Keyword string

func NewKeyword(value string) *Keyword {
	k := Keyword(value)
	return &k
}

func (s *Keyword) Accept(visitor Visitor) error {
	return visitor.VisitKeyword(s)
}

type Boolean string

func NewBoolean(value string) *Boolean {
	b := Boolean(value)
	return &b
}

func (b *Boolean) Accept(visitor Visitor) error {
	return visitor.VisitBoolean(b)
}
