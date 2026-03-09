package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
	input    textinput.Model
	err      string
	loading  bool
	width    int
	height   int
}

func newAuthModel() authModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your API key"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return authModel{
		input: ti,
	}
}

func (m authModel) SetSize(width, height int) authModel {
	m.width = width
	m.height = height
	m.input.Width = min(50, width-20)
	return m
}

func (m authModel) SetError(err string) authModel {
	m.err = err
	m.loading = false
	return m
}

func (m authModel) SetLoading(loading bool) authModel {
	m.loading = loading
	return m
}

func (m authModel) Update(msg tea.KeyMsg) (authModel, tea.Cmd) {
	if m.loading {
		return m, nil
	}

	switch msg.String() {
	case "enter":
		value := strings.TrimSpace(m.input.Value())
		if len(value) < 10 {
			m.err = "API key seems too short"
			return m, nil
		}
		m.err = ""
		m.loading = true
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m authModel) Value() string {
	return strings.TrimSpace(m.input.Value())
}

func (m authModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Welcome to Bithost"))
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	b.WriteString("Please enter your API key to get started.\n")
	b.WriteString(subtleStyle.Render("Get your key at: https://dashboard.bithost.io/api/keys"))
	b.WriteString("\n\n")

	// Input
	b.WriteString(labelStyle.Render("API Key:"))
	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString("\n")

	// Error
	if m.err != "" {
		b.WriteString("\n")
		b.WriteString(failedStyle.Render("Error: " + m.err))
		b.WriteString("\n")
	}

	// Loading
	if m.loading {
		b.WriteString("\n")
		b.WriteString(pendingStyle.Render("Authenticating..."))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter: submit • ctrl+c: quit"))

	return b.String()
}
