package repl

import (
	"raiton/evaluator"
	"raiton/lexer"
	"raiton/parser"

	tea "github.com/charmbracelet/bubbletea"
)

func (r *repl) evaluateSource(input string) tea.Cmd {
	return func() tea.Msg {
		lex := lexer.New(input)
		par := parser.New(&lex)

		node, err := par.Parse()

		if err != nil {
			return errorMsg(err)
		}

		eval := evaluator.New(r.env)

		result, err := eval.Evaluate(node)

		if err != nil {
			return errorMsg(err)
		}

		return resultMsg(result)
	}
}
