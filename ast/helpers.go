package ast

func ScopeExpressions(expressions ...Expression) *Scope {
	return &Scope{
		Definitions: []*Definition{},
		Expressions: expressions,
	}
}

func NewIdentifier(value string) *Identifier {
	ident := Identifier(value)
	return &ident
}

func NewSelector(identifiers ...*SelectorItem) *Selector {
	return &Selector{
		Items: identifiers,
	}
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

func NewApplication(arguments ...Expression) *Application {
	return &Application{
		Arguments: arguments,
	}
}

func NewArray(size uint64, elements ...Expression) *Array {
	return &Array{
		Size:     &size,
		Elements: elements,
	}
}

func NewList(elements ...Expression) *List {
	return &List{
		Elements: elements,
	}
}

func NewInteger(value int64) *Integer {
	i := Integer(value)
	return &i
}

func NewFloat(value float64) *Float {
	f := Float(value)
	return &f
}

func NewString(value string) *String {
	s := String(value)
	return &s
}

func NewKeyword(value string) *Keyword {
	k := Keyword(value)
	return &k
}

func NewBoolean(value string) *Boolean {
	b := Boolean(value)
	return &b
}
