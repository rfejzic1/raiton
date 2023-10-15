package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

type repl struct {
	textInput textinput.Model
}

func initialModel() *repl {
	ti := textinput.New()

	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 256

	return &repl{
		textInput: ti,
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
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m *repl) View() string {
	return fmt.Sprintf("%s\n(type 'exit' or ctrl+c to quit)", m.textInput.View())
}

func runRepl(ctx *cli.Context) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
