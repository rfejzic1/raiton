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

type Application struct {
	Arguments []Expression
}

func (n *Application) Type() string {
	return APPLICATION
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

type List struct {
	Elements []Expression
}

func (n *List) Type() string {
	return LIST
}

type Integer int64

func (n *Integer) Type() string {
	return INTEGER
}

type Float float64

func (n *Float) Type() string {
	return FLOAT
}

type String string

func (n *String) Type() string {
	return STRING
}

type Keyword string

func (n *Keyword) Type() string {
	return KEYWORD
}

type Boolean string

func (n *Boolean) Type() string {
	return BOOLEAN
}
