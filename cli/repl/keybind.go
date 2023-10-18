package repl

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type keyMap struct {
	Evaluate key.Binding
	Quit     key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
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
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "down"),
		),
	}
}
