package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var m = initialModel()
	if checkCommand("git") {
		m.gitStatus = greenStyle.Render("✓")
	} else {
		m.gitStatus = redStyle.Render("✗")
	}
	if checkCommand("git-lfs") {
		m.gitStatus += greenStyle.Render(" (LFS: ✓)")
	} else {
		m.gitStatus += redStyle.Render(" (LFS: ✗)")
	}

	if checkCommand("rustc") {
		m.rustStatus = greenStyle.Render("✓")
	} else {
		m.rustStatus = redStyle.Render("✗")
	}
	if checkCommand("cargo") {
		m.rustStatus += greenStyle.Render(" (Cargo: ✓)")
	} else {
		m.rustStatus += redStyle.Render(" (Cargo: ✗)")
	}
	if checkCommand("rustup") {
		m.rustStatus += greenStyle.Render(" (Rustup: ✓)")
	} else {
		m.rustStatus += redStyle.Render(" (Rustup: ✗)")
	}
	if rustVersionOK("1.88.0") {
		m.rustStatus += greenStyle.Render(" (Rust version: ✓)")
	} else {
		m.rustStatus += redStyle.Render(" (Rust version: ✗ < 1.88.0)")
	}

	if checkCommand("godot") {
		m.godotStatus = greenStyle.Render("✓")
	} else {
		m.godotStatus = orangeStyle.Render("✗")
	}

	if checkCommand("zig") {
		m.zigStatus = greenStyle.Render("✓")
	} else {
		m.zigStatus = redStyle.Render("✗")
	}

	if checkCommand("docker") {
		m.dockerStatus = greenStyle.Render("✓")
	} else {
		m.dockerStatus = redStyle.Render("✗")
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

const (
	dotChar = " • "
)

var (
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	orangeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

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

	gitStatus    string
	rustStatus   string
	godotStatus  string
	zigStatus    string
	dockerStatus string
}

type (
	frameMsg struct{}
)

func initialModel() model {
	return model{
		quitting: false,
		choice:   0,
		chosen:   false,

		gitStatus:    "",
		rustStatus:   "",
		godotStatus:  "",
		zigStatus:    "",
		dockerStatus: "",
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
	case "ctrl+c", "q", "esc":
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

	case "esc", "q":
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
			s = keywordStyle.Render("You chose Rust!")
		case 2:
			s = keywordStyle.Render("You chose Godot!")
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

func checkCommand(cmd string) bool {
	_, err := exec.Command(cmd, "--version").Output()
	if err != nil {
		_, err = exec.Command(cmd, "version").Output()
		if err != nil {
			return false
		}
		return true
	}
	return true
}

func rustVersionOK(minVersion string) bool {
	out, err := exec.Command("rustc", "--version").Output()
	if err != nil {
		return false
	}
	parts := strings.Fields(string(out))
	if len(parts) < 2 {
		return false
	}
	installed := parts[1]
	return compareVersions(installed, minVersion) >= 0
}

// compareVersions returns 1 if v1 > v2, 0 if equal, -1 if v1 < v2
func compareVersions(v1, v2 string) int {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")
	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(s1[i])
		n2, _ := strconv.Atoi(s2[i])
		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}
	return 0
}

func choicesView(m model) string {
	c := m.choice

	tpl := "Where would you like to begin?\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n\t%s\n\t%s\n\n%s\n\t%s\n\t%s\n\n%s\n\t%s\n\t%s\n\n%s\n\t%s\n\t%s\n\n%s\n\t%s\n\t%s\n\n",
		checkbox("Git / Git LFS", c == 0),
		"required | "+m.gitStatus,
		"Used for version control and large file storage",
		checkbox("Rust", c == 1),
		"required | "+m.rustStatus,
		"Used for building most of the backend and logic for Reia",
		checkbox("Godot (wip)", c == 2),
		"optional | "+m.godotStatus,
		"Godot CLI -- you may have the engine but not the CLI",
		checkbox("Zig (wip)", c == 3),
		"optional | "+m.zigStatus,
		"Zig is used for some low-level tasks but isn't used yet",
		checkbox("Docker (wip)", c == 4),
		"optional | "+m.dockerStatus,
		"Used for containerization and deployment of Reia services",
	)

	return fmt.Sprintf(tpl, choices)
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
