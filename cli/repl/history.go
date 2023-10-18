package repl

import "strings"

type history struct {
	lines []string
	index int
}

func newHistory() history {
	return history{
		lines: []string{},
		index: -1,
	}
}

func (h *history) add(line string) {
	if strings.TrimSpace(line) != "" {
		h.lines = append(h.lines, line)
	}
}

func (h *history) previous() string {
	if h.index < len(h.lines)-1 {
		h.index += 1
	}

	if h.index >= 0 && h.index < len(h.lines) {
		line := h.lines[len(h.lines)-1-h.index]
		return line
	}

	return ""
}

func (h *history) next() string {
	if h.index >= 0 {
		h.index -= 1
	}

	if h.index >= 0 && h.index < len(h.lines) {
		line := h.lines[len(h.lines)-1-h.index]
		return line
	}

	return ""
}

func (h *history) reset() {
	h.index = -1
}
