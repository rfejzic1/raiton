package ast

const (
	SCOPE         = "scope"
	DEFINITION    = "definition"
	APPLICATION   = "application"
	IDENTIFIER    = "identifier"
	SELECTOR      = "selector"
	SELECTOR_ITEM = "selector_item"
	FUNCTION      = "function"
	CONDITIONAL   = "conditional"
	RECORD        = "record"
	ARRAY         = "array"
	LIST          = "list"
	INTEGER       = "integer"
	FLOAT         = "float"
	STRING        = "string"
	KEYWORD       = "keyword"
	BOOLEAN       = "boolean"
)

type Node interface {
	Type() string
}

type Scope struct {
	Definitions []*Definition
	Expressions []Expression
}

func (n *Scope) Type() string {
	return SCOPE
}

func ScopeExpressions(expressions ...Expression) *Scope {
	return &Scope{
		Definitions: []*Definition{},
		Expressions: expressions,
	}
}

type Definition struct {
	Identifier Identifier
	Expression Expression
}

func (n *Definition) Type() string {
	return DEFINITION
}

// *** Expressions ***

type Expression interface {
	Node
}

type Identifier string

func (n *Identifier) Type() string {
	return IDENTIFIER
}

func NewIdentifier(value string) *Identifier {
	ident := Identifier(value)
	return &ident
}

type Selector struct {
	Items []*SelectorItem
}

func (n *Selector) Type() string {
	return SELECTOR
}

type SelectorItem struct {
	Identifier *Identifier
	Index      *Integer
}

func (n *SelectorItem) Type() string {
	return SELECTOR_ITEM
}

func NewIdentifierSelector(ident *Identifier) *SelectorItem {
	return &SelectorItem{
		Identifier: ident,
	}
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

type Application struct {
	Arguments []Expression
}

func (n *Application) Type() string {
	return APPLICATION
}

func NewApplication(arguments ...Expression) *Application {
	return &Application{
		Arguments: arguments,
	}
}

type Function struct {
	Parameters []*Identifier
	Body       *Scope
}

func (n *Function) Type() string {
	return FUNCTION
}

type Conditional struct {
	Condition   Expression
	Consequence *Scope
	Alternative *Scope
}

func (n *Conditional) Type() string {
	return CONDITIONAL
}

type Record struct {
	Fields map[Identifier]Expression
}

func (n *Record) Type() string {
	return RECORD
}

type Array struct {
	Size     *uint64
	Elements []Expression
}

func (n *Array) Type() string {
	return ARRAY
}

func NewArray(size uint64, elements ...Expression) *Array {
	return &Array{
		Size:     &size,
		Elements: elements,
	}
}

type List struct {
	Elements []Expression
}

func (n *List) Type() string {
	return LIST
}

func NewList(elements ...Expression) *List {
	return &List{
		Elements: elements,
	}
}

type Integer int64

func (n *Integer) Type() string {
	return INTEGER
}

func NewInteger(value int64) *Integer {
	i := Integer(value)
	return &i
}

type Float float64

func (n *Float) Type() string {
	return FLOAT
}

func NewFloat(value float64) *Float {
	f := Float(value)
	return &f
}

type String string

func (n *String) Type() string {
	return STRING
}

func NewString(value string) *String {
	s := String(value)
	return &s
}

type Keyword string

func (n *Keyword) Type() string {
	return KEYWORD
}

func NewKeyword(value string) *Keyword {
	k := Keyword(value)
	return &k
}

type Boolean string

func (n *Boolean) Type() string {
	return BOOLEAN
}

func NewBoolean(value string) *Boolean {
	b := Boolean(value)
	return &b
}
