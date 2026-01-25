package tui

// DashboardHelp returns a short help string for the dashboard
func DashboardHelp() string {
	return helpStyle.Render("n: new server • s: ssh keys • t: top up • r: refresh • enter: details • d: delete • q: quit")
}

// DetailHelp returns help for server detail view
func DetailHelp() string {
	return helpStyle.Render("c: copy ssh • d: delete • q/esc: back")
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
