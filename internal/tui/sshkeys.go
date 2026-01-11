package tui

import (
	"os"
	"strings"

	"github.com/bithostio/bh/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sshkeysModel struct {
	keys    []api.SSHKey
	cursor  int
	width   int
	height  int
	loading bool

	// Add key form
	adding     bool
	addStep    int // 0 = label, 1 = key/path
	labelInput textinput.Model
	keyInput   textinput.Model
	addError   string
}

const (
	addStepLabel = iota
	addStepKey
)

func newSSHKeysModel() sshkeysModel {
	labelInput := textinput.New()
	labelInput.Placeholder = "my-laptop"
	labelInput.CharLimit = 64
	labelInput.Width = 40

	keyInput := textinput.New()
	keyInput.Placeholder = "~/.ssh/id_rsa.pub or paste key"
	keyInput.CharLimit = 2048
	keyInput.Width = 50

	return sshkeysModel{
		loading:    true,
		labelInput: labelInput,
		keyInput:   keyInput,
	}
}

func (m sshkeysModel) SetKeys(keys []api.SSHKey) sshkeysModel {
	m.keys = keys
	m.loading = false
	return m
}

func (m sshkeysModel) SetSize(width, height int) sshkeysModel {
	m.width = width
	m.height = height
	m.labelInput.Width = min(40, width-20)
	m.keyInput.Width = min(50, width-20)
	return m
}

func (m sshkeysModel) StartAdding() sshkeysModel {
	m.adding = true
	m.addStep = addStepLabel
	m.addError = ""
	m.labelInput.Reset()
	m.keyInput.Reset()
	m.labelInput.Focus()
	return m
}

func (m sshkeysModel) CancelAdding() sshkeysModel {
	m.adding = false
	m.addStep = addStepLabel
	m.addError = ""
	m.labelInput.Blur()
	m.keyInput.Blur()
	return m
}

func (m sshkeysModel) Update(msg tea.KeyMsg) (sshkeysModel, tea.Cmd) {
	if m.adding {
		return m.updateAddForm(msg)
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.keys)-1 {
			m.cursor++
		}
	}
	return m, nil
}

func (m sshkeysModel) updateAddForm(msg tea.KeyMsg) (sshkeysModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m = m.CancelAdding()
		return m, nil
	case "enter":
		switch m.addStep {
		case addStepLabel:
			label := strings.TrimSpace(m.labelInput.Value())
			if label == "" {
				m.addError = "Label cannot be empty"
				return m, nil
			}
			m.addError = ""
			m.addStep = addStepKey
			m.labelInput.Blur()
			m.keyInput.Focus()
			return m, nil
		case addStepKey:
			// Will be handled by tui.go to call the API
			return m, nil
		}
	default:
		var cmd tea.Cmd
		if m.addStep == addStepLabel {
			m.labelInput, cmd = m.labelInput.Update(msg)
		} else {
			m.keyInput, cmd = m.keyInput.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

// GetKeyToAdd returns the label and key content for adding
// Returns empty strings if not ready
func (m sshkeysModel) GetKeyToAdd() (label, key string, ok bool) {
	if !m.adding || m.addStep != addStepKey {
		return "", "", false
	}

	label = strings.TrimSpace(m.labelInput.Value())
	keyOrPath := strings.TrimSpace(m.keyInput.Value())

	if label == "" || keyOrPath == "" {
		return "", "", false
	}

	// Check if it's a file path
	if strings.HasPrefix(keyOrPath, "~") || strings.HasPrefix(keyOrPath, "/") || strings.HasPrefix(keyOrPath, ".") {
		// Expand ~ to home directory
		if strings.HasPrefix(keyOrPath, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				keyOrPath = strings.Replace(keyOrPath, "~", home, 1)
			}
		}
		// Try to read file
		data, err := os.ReadFile(keyOrPath)
		if err != nil {
			return "", "", false
		}
		key = strings.TrimSpace(string(data))
	} else {
		// Assume it's the key content
		key = keyOrPath
	}

	return label, key, true
}

func (m sshkeysModel) SetAddError(err string) sshkeysModel {
	m.addError = err
	return m
}

func (m sshkeysModel) View() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("SSH KEYS")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	// Show add form if adding
	if m.adding {
		b.WriteString(m.renderAddForm())
		return b.String()
	}

	if m.loading {
		b.WriteString(subtleStyle.Render("Loading..."))
		b.WriteString("\n\n")
		b.WriteString(SSHKeysHelp())
		return b.String()
	}

	if len(m.keys) == 0 {
		b.WriteString(subtleStyle.Render("No SSH keys found. Press 'a' to add one."))
		b.WriteString("\n\n")
		b.WriteString(SSHKeysHelp())
		return b.String()
	}

	for i, key := range m.keys {
		keyPreview := key.Key
		if len(keyPreview) > 40 {
			keyPreview = keyPreview[:37] + "..."
		}

		label := padRight(key.Label, 20)
		preview := subtleStyle.Render(keyPreview)

		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + selectedStyle.Render(label) + "  " + preview + "\n")
		} else {
			b.WriteString("  " + label + "  " + preview + "\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(SSHKeysHelp())

	return b.String()
}

func (m sshkeysModel) renderAddForm() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ADD SSH KEY"))
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", 40)))
	b.WriteString("\n\n")

	// Label input
	labelLabel := "Label:"
	if m.addStep == addStepLabel {
		labelLabel = selectedStyle.Render("> ") + labelLabel
	} else {
		labelLabel = "  " + labelLabel
	}
	b.WriteString(labelLabel)
	b.WriteString("\n  ")
	if m.addStep == addStepLabel {
		b.WriteString(m.labelInput.View())
	} else {
		b.WriteString(valueStyle.Render(m.labelInput.Value()))
	}
	b.WriteString("\n\n")

	// Key input (only show if past label step)
	if m.addStep >= addStepKey {
		b.WriteString(selectedStyle.Render("> ") + "Public Key (file path or paste):")
		b.WriteString("\n  ")
		b.WriteString(m.keyInput.View())
		b.WriteString("\n")
	}

	// Error message
	if m.addError != "" {
		b.WriteString("\n")
		b.WriteString(failedStyle.Render("Error: " + m.addError))
		b.WriteString("\n")
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter: continue • esc: cancel"))

	return b.String()
}
