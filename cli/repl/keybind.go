package repl

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Evaluate key.Binding
	Quit     key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up", "previous line"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("down", "next line"),
		),
		Evaluate: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "evaluate"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func viewportKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		Up: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "down"),
		),
	}
}
