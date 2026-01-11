package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Adaptive colors for light/dark terminal support
var (
	primaryColor = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	successColor = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	warningColor = lipgloss.AdaptiveColor{Light: "#FF9900", Dark: "#FFAA00"}
	errorColor   = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF5555"}
	subtleColor  = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#626262"}
	mutedColor   = lipgloss.AdaptiveColor{Light: "#626262", Dark: "#9B9B9B"}
)

// Base styles
var (
	// Title style for headers
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	// Status styles
	activeStyle = lipgloss.NewStyle().
			Foreground(successColor)

	pendingStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	failedStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Subtle text
	subtleStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Help bar style
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Modal styles
	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(50)

	modalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(warningColor).
			MarginBottom(1)

	// Value styles for detail view
	labelStyle = lipgloss.NewStyle().
			Width(14).
			Foreground(subtleColor)

	valueStyle = lipgloss.NewStyle().
			Bold(true)

	// Wizard step indicator
	stepStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

// FormatStatus returns a styled status string
func FormatStatus(status string, pending, power bool) string {
	if pending {
		return pendingStyle.Render("pending")
	}
	switch status {
	case "active":
		if power {
			return activeStyle.Render("active")
		}
		return pendingStyle.Render("powered off")
	case "failed":
		return failedStyle.Render("failed")
	default:
		return status
	}
}

// FormatMoney formats a price
func FormatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// FormatEnabled returns styled enabled/disabled
func FormatEnabled(enabled bool) string {
	if enabled {
		return activeStyle.Render("Enabled")
	}
	return mutedStyle.Render("Disabled")
}
