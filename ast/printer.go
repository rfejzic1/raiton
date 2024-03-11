package ast

import (
	"fmt"
	"strings"
)

type Printer struct {
	node Node
	sb   strings.Builder
}

func NewPrinter(node Node) *Printer {
	return &Printer{
		node: node,
	}
}

func (p *Printer) String() string {
	p.print(p.node)
	return p.sb.String()
}

func (p *Printer) write(s string) {
	p.sb.WriteString(s)
}

func (p *Printer) writeln() {
	p.write("\n")
}

func (e *Printer) print(node Node) error {
	switch n := node.(type) {
	case *Scope:
		return e.scope(n)
	case *Definition:
		return e.definition(n)
	case *Identifier:
		return e.identifier(n)
	case *Selector:
		return e.selector(n)
	case *SelectorItem:
		return e.selectorItem(n)
	case *Application:
		return e.application(n)
	case *Function:
		return e.function(n)
	case *Conditional:
		return e.conditional(n)
	case *Record:
		return e.record(n)
	case *Array:
		return e.array(n)
	case *List:
		return e.list(n)
	case *String:
		return e.string(n)
	case *Integer:
		return e.integer(n)
	case *Float:
		return e.float(n)
	case *Keyword:
		return e.keyword(n)
	case *Boolean:
		return e.boolean(n)
	default:
		panic("unhandled ast type")
	}
}

/*** Print Methods ***/

func (p *Printer) scope(n *Scope) error {
	for i, d := range n.Definitions {
		if err := p.print(d); err != nil {
			return err
		}
		if i != len(n.Definitions)-1 {
			p.write(" ")
		}
	}

	if len(n.Definitions) > 0 && len(n.Expressions) > 0 {
		p.writeln()
	}

	for i, e := range n.Expressions {
		if err := p.print(e); err != nil {
			return err
		}
		if i != len(n.Expressions)-1 {
			p.write(" ")
		}
	}

	return nil
}

func (p *Printer) definition(n *Definition) error {
	p.write(string(n.Identifier))

	_, is_scope := n.Expression.(*Scope)

	if is_scope {
		p.write(" { ")

		if err := p.print(n.Expression); err != nil {
			return err
		}

		p.write(" }")
		p.writeln()
	} else {
		p.write(": ")
		if err := p.print(n.Expression); err != nil {
			return err
		}
		p.writeln()
	}

	return nil
}

func (p *Printer) identifier(n *Identifier) error {
	p.write(string(*n))
	return nil
}

func (p *Printer) selector(n *Selector) error {
	l := len(n.Items)

	for i, ip := range n.Items {
		if err := p.print(ip); err != nil {
			return nil
		}

		if i != l-1 {
			p.write(".")
		}

	}

	return nil
}

func (p *Printer) selectorItem(n *SelectorItem) error {
	if n.Identifier != nil {
		p.write(string(*n.Identifier))
	} else {
		p.write(fmt.Sprintf("%d", *n.Index))
	}

	return nil
}

func (p *Printer) application(n *Application) error {
	p.write("(")

	for _, a := range n.Arguments {
		if err := p.print(a); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write(")")

	return nil
}

func (p *Printer) conditional(n *Conditional) error {
	p.write("if ")

	if err := p.print(n.Condition); err != nil {
		return err
	}

	if err := p.print(n.Consequence); err != nil {
		return err
	}

	p.write(" else ")

	if err := p.print(n.Alternative); err != nil {
		return err
	}

	return nil
}

func (p *Printer) function(n *Function) error {
	p.write("\\")

	for _, param := range n.Parameters {
		p.write(string(*param))
		p.write(" ")
	}

	p.write("{ ")

	if err := p.print(n.Body); err != nil {
		return err
	}

	p.write(" }")

	return nil
}

func (p *Printer) record(n *Record) error {
	p.write("{ ")

	for ident, expr := range n.Fields {
		p.write(string(ident))
		p.write(": ")

		if err := p.print(expr); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("}")

	return nil
}

func (p *Printer) array(n *Array) error {
	p.write("[ ")

	for _, expr := range n.Elements {
		if err := p.print(expr); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("]")

	return nil
}

func (p *Printer) list(n *List) error {
	p.write("[ ")

	for _, expr := range n.Elements {
		if err := p.print(expr); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("]")

	return nil
}

func (p *Printer) integer(n *Integer) error {
	p.write(fmt.Sprintf("%d", *n))
	return nil
}

func (p *Printer) float(n *Float) error {
	p.write(fmt.Sprintf("%g", *n))
	return nil
}

func (p *Printer) string(n *String) error {
	p.write(fmt.Sprintf("\"%s\"", *n))
	return nil
}

func (p *Printer) keyword(n *Keyword) error {
	p.write(string(*n))
	return nil
}

func (p *Printer) boolean(n *Boolean) error {
	p.write(string(*n))
	return nil
}
