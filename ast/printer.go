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
	p.node.Accept(p)
	return p.sb.String()
}

func (p *Printer) write(s string) {
	p.sb.WriteString(s)
}

func (p *Printer) writeln() {
	p.write("\n")
}

/*** Visitor Methods ***/

func (p *Printer) VisitScope(n *Scope) error {
	for i, d := range n.Definitions {
		if err := d.Accept(p); err != nil {
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
		if err := e.Accept(p); err != nil {
			return err
		}
		if i != len(n.Expressions)-1 {
			p.write(" ")
		}
	}

	return nil
}

func (p *Printer) VisitDefinition(n *Definition) error {
	p.write(string(n.Identifier))

	_, is_scope := n.Expression.(*Scope)

	if is_scope {
		p.write(" { ")

		if err := n.Expression.Accept(p); err != nil {
			return err
		}

		p.write(" }")
		p.writeln()
	} else {
		p.write(": ")
		if err := n.Expression.Accept(p); err != nil {
			return err
		}
		p.writeln()
	}

	return nil
}

func (p *Printer) VisitIdentifier(n *Identifier) error {
	p.write(string(*n))
	return nil
}

func (p *Printer) VisitSelector(n *Selector) error {
	l := len(n.Items)

	for i, ip := range n.Items {
		if err := ip.Accept(p); err != nil {
			return nil
		}

		if i != l-1 {
			p.write(".")
		}

	}

	return nil
}

func (p *Printer) VisitSelectorItem(n *SelectorItem) error {
	if n.Identifier != nil {
		p.write(string(*n.Identifier))
	} else {
		p.write(fmt.Sprintf("%d", *n.Index))
	}

	return nil
}

func (p *Printer) VisitApplication(n *Application) error {
	p.write("(")

	for _, a := range n.Arguments {
		if err := a.Accept(p); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write(")")

	return nil
}

func (p *Printer) VisitFunction(n *FunctionLiteral) error {
	p.write("\\")

	for _, param := range n.Parameters {
		p.write(string(*param))
		p.write(" ")
	}

	p.write("{ ")

	if err := n.Body.Accept(p); err != nil {
		return err
	}

	p.write(" }")

	return nil
}

func (p *Printer) VisitRecord(n *RecordLiteral) error {
	p.write("{ ")

	for ident, expr := range n.Fields {
		p.write(string(ident))
		p.write(": ")

		if err := expr.Accept(p); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("}")

	return nil
}

func (p *Printer) VisitArray(n *ArrayLiteral) error {
	p.write("[ ")

	for _, expr := range n.Elements {
		if err := expr.Accept(p); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("]")

	return nil
}

func (p *Printer) VisitSlice(n *SliceLiteral) error {
	p.write("[ ")

	for _, expr := range n.Elements {
		if err := expr.Accept(p); err != nil {
			return err
		}

		p.write(" ")
	}

	p.write("]")

	return nil
}

func (p *Printer) VisitInteger(n *IntegerLiteral) error {
	p.write(fmt.Sprintf("%d", *n))
	return nil
}

func (p *Printer) VisitFloat(n *FloatLiteral) error {
	p.write(fmt.Sprintf("%g", *n))
	return nil
}

func (p *Printer) VisitString(n *StringLiteral) error {
	p.write(string(*n))
	return nil
}

func (p *Printer) VisitBoolean(n *BooleanLiteral) error {
	p.write(string(*n))
	return nil
}
