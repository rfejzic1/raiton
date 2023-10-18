package repl

import (
	"fmt"
	"raiton/object"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

type repl struct {
	loading   bool
	width     int
	height    int
	keys      keyMap
	lines     []string
	viewport  viewport.Model
	textInput textinput.Model
	history   history
	env       *object.Environment
}

type errorMsg error
type resultMsg object.Object

var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(105))
var promptStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(105))
var dimmedStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
var expressionStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(7))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(9))

func initialModel() *repl {
	vp := viewport.New(0, 0)
	vp.KeyMap = viewportKeyMap()

	ti := textinput.New()
	ti.Focus()
	ti.PromptStyle = promptStyle

	lines := []string{"Raiton " + titleStyle.Render("v0.0.1")}
	vp.SetContent(strings.Join(lines, "\n"))

	return &repl{
		loading:   true,
		keys:      defaultKeyMap(),
		lines:     lines,
		viewport:  vp,
		textInput: ti,
		env:       object.NewEnvironment(),
		history:   newHistory(),
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
		case key.Matches(msg, m.keys.Up):
			return m.previousItem(msg)
		case key.Matches(msg, m.keys.Down):
			return m.nextItem(msg)
		case key.Matches(msg, m.keys.Evaluate):
			return m.evaluate(msg)
		}
	case resultMsg:
		m.addLine(msg.Inspect())
		return m, nil
	case errorMsg:
		m.addLine(errorStyle.Render(msg.Error()))
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *repl) evaluate(msg tea.Msg) (*repl, tea.Cmd) {
	rawInput := m.textInput.Value()
	input := strings.TrimSpace(rawInput)

	m.textInput.Reset()

	if input == "exit" {
		return m, tea.Quit
	}

	line := fmt.Sprintf("%s %s", promptStyle.Render(">"), rawInput)
	m.addLine(line)

	if input == "" {
		return m, nil
	}

	m.history.add(input)

	return m, m.evaluateSource(input)
}

func (m *repl) previousItem(msg tea.Msg) (*repl, tea.Cmd) {
	line := m.history.previous()
	m.textInput.SetValue(line)
	return m, nil
}

func (m *repl) nextItem(msg tea.Msg) (*repl, tea.Cmd) {
	line := m.history.next()
	m.textInput.SetValue(line)
	return m, nil
}

func (r *repl) addLine(line string) {
	r.lines = append(r.lines, expressionStyle.Render(line))
	r.history.reset()
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
	s.WriteString(dimmedStyle.Render("(type 'exit' or ctrl+c to quit)"))

	return fmt.Sprintf(s.String())
}

func Run(ctx *cli.Context) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
