package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

const (
	dotChar = " • "
)

var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

type model struct {
	quitting bool
	choice   int
	chosen   bool
}

type (
	frameMsg struct{}
)

func initialModel() model {
	return model{
		quitting: false,
		choice:   0,
		chosen:   false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.chosen {
			return m.handleMainScreen(msg)
		} else {
			return m.handleSubScreen(msg)
		}
	}
	return m, nil
}

func (m model) handleMainScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.quitting = true
		return m, tea.Quit
	case "j", "down":
		m.choice++
		if m.choice > 4 {
			m.choice = 4
		}
	case "k", "up":
		m.choice--
		if m.choice < 0 {
			m.choice = 0
		}
	case "enter":
		m.chosen = true
		return m, frame()
	}

	return m, nil
}

func (m model) handleSubScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "esc":
		m.chosen = false

	}
	return m, nil
}

func (m model) View() string {
	var s string
	if m.quitting {
		return "\n  See you later!\n\n"
	}
	if !m.chosen {
		s = choicesView(m)
	} else {
		switch m.choice {
		case 0:
			s = keywordStyle.Render("You chose Git / Git LFS!")
		case 1:
			s = keywordStyle.Render("You chose Godot!")
		case 2:
			s = keywordStyle.Render("You chose Rust!")
		case 3:
			s = keywordStyle.Render("You chose Zig!")
		case 4:
			s = keywordStyle.Render("You chose Docker!")
		default:
			s = keywordStyle.Render("Unknown choice")
		}
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

func choicesView(m model) string {
	c := m.choice

	tpl := "Where would you like to begin?\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("ctrl+c, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		checkbox("Git / Git LFS", c == 0),
		checkbox("Godot", c == 1),
		checkbox("Rust", c == 2),
		checkbox("Zig", c == 3),
		checkbox("Docker", c == 4),
	)

	return fmt.Sprintf(tpl, choices)
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
