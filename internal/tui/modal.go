package tui

import (
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type modalType int

const (
	modalDelete modalType = iota
	modalTopup
)

type modalModel struct {
	active    bool
	modalType modalType
	message   string
	serverID  int
	confirmed bool
}

func newModalModel(message string, serverID int) modalModel {
	return modalModel{
		active:    true,
		modalType: modalDelete,
		message:   message,
		serverID:  serverID,
	}
}

func newTopupModal() modalModel {
	return modalModel{
		active:    true,
		modalType: modalTopup,
	}
}

func (m modalModel) Update(msg tea.KeyMsg) (modalModel, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.confirmed = true
		m.active = false
	case "n", "N", "esc", "q":
		m.confirmed = false
		m.active = false
	}
	return m, nil
}

func (m modalModel) View(background string, width, height int) string {
	if !m.active {
		return background
	}

	var content strings.Builder

	switch m.modalType {
	case modalDelete:
		content.WriteString(modalTitleStyle.Render("Confirm Delete"))
		content.WriteString("\n\n")
		content.WriteString(m.message)
		content.WriteString("\n\n")
		content.WriteString(ModalHelp())
	case modalTopup:
		content.WriteString(titleStyle.Render("Top Up Balance"))
		content.WriteString("\n\n")
		content.WriteString("To add funds to your account, visit:\n\n")
		content.WriteString(selectedStyle.Render("  https://dashboard.bithost.io"))
		content.WriteString("\n\n")
		content.WriteString("Open in browser?")
		content.WriteString("\n\n")
		content.WriteString(helpStyle.Render("y: open browser • q/n/esc: cancel"))
	}

	modal := modalStyle.Render(content.String())

	// Calculate position to center the modal
	modalWidth := lipgloss.Width(modal)
	modalHeight := lipgloss.Height(modal)

	x := max((width-modalWidth)/2, 0)
	y := max((height-modalHeight)/2, 0)

	return placeOverlay(x, y, modal, background, width, height)
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// placeOverlay places a modal overlay on top of the background
func placeOverlay(x, y int, overlay, background string, width, height int) string {
	bgLines := strings.Split(background, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// Ensure we have enough lines
	for len(bgLines) < height {
		bgLines = append(bgLines, "")
	}

	// Overlay each line
	for i, overlayLine := range overlayLines {
		bgY := y + i
		if bgY >= 0 && bgY < len(bgLines) {
			bgLines[bgY] = insertAt(bgLines[bgY], x, overlayLine, width)
		}
	}

	return strings.Join(bgLines, "\n")
}

// insertAt inserts overlay text into a background line at position x
func insertAt(bgLine string, x int, overlay string, _ int) string {
	bgRunes := []rune(bgLine)

	// Pad background line to reach x position
	for len(bgRunes) < x {
		bgRunes = append(bgRunes, ' ')
	}

	overlayRunes := []rune(overlay)
	overlayWidth := lipgloss.Width(overlay)

	// Build: background prefix + overlay + background suffix
	result := append(bgRunes[:x], overlayRunes...)

	afterX := x + overlayWidth
	if afterX < len(bgRunes) {
		result = append(result, bgRunes[afterX:]...)
	}

	return string(result)
}
