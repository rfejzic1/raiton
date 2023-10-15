package cli

import (
	"fmt"
	"raiton/evaluator"
	"raiton/lexer"
	"raiton/object"
	"raiton/parser"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

type repl struct {
	textInput textinput.Model
	env       *object.Environment
	lines     []string
}

type errorMsg error
type resultMsg object.Object

func initialModel() *repl {
	ti := textinput.New()

	ti.Focus()
	ti.CharLimit = 256

	return &repl{
		textInput: ti,
		env:       object.NewEnvironment(),
		lines:     []string{},
	}
}

func (m *repl) Init() tea.Cmd {
	return textinput.Blink
}

func (m *repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			rawInput := m.textInput.Value()
			input := strings.TrimSpace(rawInput)

			m.textInput.Reset()

			if input == "exit" {
				return m, tea.Quit
			}

			line := fmt.Sprintf("> %s", rawInput)
			m.lines = append(m.lines, line)

			if input == "" {
				return m, nil
			}

			return m, m.evaluateSource(input)
		}
	case resultMsg:
		m.lines = append(m.lines, msg.Inspect())
		return m, nil
	case errorMsg:
		m.lines = append(m.lines, msg.Error())
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

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

func (m *repl) View() string {
	var s strings.Builder

	s.WriteString(lipgloss.JoinVertical(lipgloss.Left, m.lines...))

	if len(m.lines) != 0 {
		s.WriteString("\n")
	}

	s.WriteString(m.textInput.View())
	s.WriteString("\n")
	s.WriteString("(type 'exit' or ctrl+c to quit)")

	return fmt.Sprintf(s.String())
}

func runRepl(ctx *cli.Context) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
