package object

import (
	"fmt"
	"strings"

	"github.com/rfejzic1/raiton/ast"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	BOOLEAN   = "boolean"
	CHARACTER = "character"
	INTEGER   = "integer"
	FLOAT     = "float"
	STRING    = "string"
	ARRAY     = "array"
	SLICE     = "slice"
	RECORD    = "record"
	FUNCTION  = "function"
	BUILTIN   = "builtin"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN }

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func BoxBoolean(value bool) *Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

type Character struct {
	Value string
}

func (c *Character) Inspect() string { return fmt.Sprintf("'%s'", c.Value) }

func (c *Character) Type() ObjectType { return CHARACTER }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Type() ObjectType { return INTEGER }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string { return fmt.Sprintf("%g", f.Value) }

func (f *Float) Type() ObjectType { return FLOAT }

type String struct {
	Value string
}

func (s *String) Inspect() string { return fmt.Sprintf(`"%s"`, s.Value) }

func (s *String) Type() ObjectType { return STRING }

type Array struct {
	Value []Object
	Size  uint64
}

func (a *Array) Inspect() string {
	strs := []string{}

	for _, o := range a.Value {
		strs = append(strs, o.Inspect())
	}

	return fmt.Sprintf("[%d: %s]", a.Size, strings.Join(strs, " "))
}

func (a *Array) Type() ObjectType { return ARRAY }

type Slice struct {
	Value *Array
}

func (s *Slice) Inspect() string {
	strs := []string{}

	for _, o := range s.Value.Value {
		strs = append(strs, o.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, " "))
}

func (s *Slice) Type() ObjectType { return SLICE }

type Record struct {
	Value map[string]Object
}

func (r *Record) Inspect() string {
	strs := []string{}

	for field, obj := range r.Value {
		str := fmt.Sprintf("%s: %s", field, obj.Inspect())
		strs = append(strs, str)
	}

	return fmt.Sprintf("{ %s }", strings.Join(strs, " "))
}

func (r *Record) Type() ObjectType { return RECORD }

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.Scope
	Environment *Environment
}

func (f *Function) Inspect() string {
	var sb strings.Builder

	sb.WriteString("\\")

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, string(*p))
	}

	sb.WriteString(strings.Join(params, " "))

	sb.WriteString(" { ")
	p := ast.NewPrinter(f.Body)
	sb.WriteString(p.String())
	sb.WriteString(" }")

	return sb.String()
}

func (f *Function) Type() ObjectType { return FUNCTION }

type BuiltinFunction func(v ast.Visitor, args ...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFunction
}

func MakeBuiltin(fn BuiltinFunction) *Builtin {
	return &Builtin{
		Fn: fn,
	}
}

func (b *Builtin) Type() ObjectType { return BUILTIN }

func (b *Builtin) Inspect() string { return "builtin function" }
