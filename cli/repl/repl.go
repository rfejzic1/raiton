package repl

import (
	"fmt"
	"raiton/object"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

type repl struct {
	loading   bool
	width     int
	height    int
	keys      keyMap
	viewport  viewport.Model
	textInput textinput.Model
	env       *object.Environment
	lines     []string
}

type errorMsg error
type resultMsg object.Object

func initialModel() *repl {
	vp := viewport.New(0, 0)
	ti := textinput.New()

	ti.Focus()
	vp.KeyMap = viewportKeyMap()

	lines := []string{"Raiton v0.0.1"}
	vp.SetContent(strings.Join(lines, "\n"))

	return &repl{
		loading:   true,
		viewport:  vp,
		textInput: ti,
		keys:      defaultKeyMap(),
		env:       object.NewEnvironment(),
		lines:     lines,
	}
}

func (m *repl) Init() tea.Cmd {
	return textinput.Blink
}

func (m *repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd, vpCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.loading {
			m.loading = false
		}
		m.width = msg.Width
		m.height = msg.Height
		m.computeViewportHeight()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Evaluate):
			rawInput := m.textInput.Value()
			input := strings.TrimSpace(rawInput)

			m.textInput.Reset()

			if input == "exit" {
				return m, tea.Quit
			}

			line := fmt.Sprintf("> %s", rawInput)
			m.addLine(line)

			if input == "" {
				return m, nil
			}

			return m, m.evaluateSource(input)
		}
	case resultMsg:
		m.addLine(msg.Inspect())
		return m, nil
	case errorMsg:
		m.addLine(msg.Error())
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (r *repl) addLine(line string) {
	r.lines = append(r.lines, line)
	r.computeViewportHeight()
	r.viewport.SetContent(strings.Join(r.lines, "\n"))
	r.viewport.GotoBottom()
}

func (r *repl) computeViewportHeight() {
	const offset = 2
	r.viewport.Width = r.width
	r.viewport.Height = min(len(r.lines), r.height-offset)
}

func (m *repl) View() string {
	if m.loading {
		return "Loading..."
	}

	var s strings.Builder

	s.WriteString(m.viewport.View())

	if len(m.lines) > 0 {
		s.WriteString("\n")
	}

	s.WriteString(m.textInput.View())
	s.WriteString("\n")
	s.WriteString("(type 'exit' or ctrl+c to quit)")

	return fmt.Sprintf(s.String())
}

func Run(ctx *cli.Context) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
