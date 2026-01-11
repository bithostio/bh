package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the TUI
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Back    key.Binding
	New     key.Binding
	Delete  key.Binding
	Refresh key.Binding
	SSHKeys key.Binding
	Help    key.Binding
	Quit    key.Binding
	Toggle  key.Binding
	Yes     key.Binding
	No      key.Binding
}

// DefaultKeyMap returns the default keybindings
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new server"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	SSHKeys: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "ssh keys"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "tab"),
		key.WithHelp("space", "toggle"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "no"),
	),
}

// ShortHelp returns a short help string for the dashboard
func DashboardHelp() string {
	return helpStyle.Render("n: new server • s: ssh keys • t: top up • r: refresh • enter: details • d: delete • q: quit")
}

// DetailHelp returns help for server detail view
func DetailHelp() string {
	return helpStyle.Render("d: delete • q/esc: back")
}

// WizardHelp returns help for the wizard
func WizardHelp() string {
	return helpStyle.Render("↑/↓: navigate • enter: select • q/esc: back")
}

// SSHKeysHelp returns help for SSH keys view
func SSHKeysHelp() string {
	return helpStyle.Render("a: add key • q/esc: back")
}

// ModalHelp returns help for modal dialogs
func ModalHelp() string {
	return helpStyle.Render("y: confirm • q/n/esc: cancel")
}
