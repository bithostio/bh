package tui

import (
	"fmt"
	"strings"

	"github.com/bithostio/bh/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dashboardModel struct {
	servers []api.Server
	user    *api.UserResponse
	cursor  int
	width   int
	height  int
	loading bool
}

func newDashboardModel() dashboardModel {
	return dashboardModel{
		loading: true,
	}
}

func (m dashboardModel) SetServers(servers []api.Server) dashboardModel {
	// Filter to only active/pending servers
	m.servers = filterActiveServers(servers)
	m.loading = false
	if m.cursor >= len(m.servers) && len(m.servers) > 0 {
		m.cursor = len(m.servers) - 1
	}
	return m
}

func filterActiveServers(servers []api.Server) []api.Server {
	var active []api.Server
	for _, s := range servers {
		if s.Status == "active" || s.Pending {
			active = append(active, s)
		}
	}
	return active
}

func (m dashboardModel) SetUser(user *api.UserResponse) dashboardModel {
	m.user = user
	return m
}

func (m dashboardModel) SetSize(width, height int) dashboardModel {
	m.width = width
	m.height = height
	return m
}

func (m dashboardModel) selectedServer() *api.Server {
	if len(m.servers) == 0 || m.cursor >= len(m.servers) {
		return nil
	}
	return &m.servers[m.cursor]
}

func (m dashboardModel) Update(msg tea.KeyMsg) (dashboardModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.servers)-1 {
			m.cursor++
		}
	}
	return m, nil
}

func (m dashboardModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Menu bar
	b.WriteString(m.renderMenuBar())
	b.WriteString("\n\n")

	// Server list
	b.WriteString(m.renderServerList())

	// Help
	b.WriteString("\n\n")
	b.WriteString(DashboardHelp())

	return b.String()
}

func (m dashboardModel) renderHeader() string {
	title := titleStyle.Render("bithost.io")

	var right string
	if m.user != nil {
		balance := fmt.Sprintf("Balance: %s", FormatMoney(m.user.Balance))
		if m.user.ServerLimit > 0 {
			serverCount := fmt.Sprintf("Servers: %d/%d", len(m.servers), m.user.ServerLimit)
			right = subtleStyle.Render(balance + "  " + serverCount)
		} else {
			right = subtleStyle.Render(balance)
		}
	}

	// Calculate spacing (ensure at least 1 space, never negative)
	leftWidth := lipgloss.Width(title)
	rightWidth := lipgloss.Width(right)
	spacing := max(m.width-leftWidth-rightWidth-4, 1)
	header := title + strings.Repeat(" ", spacing) + right

	var userLine string
	if m.user != nil {
		userLine = subtleStyle.Render(m.user.Email)
	}

	return header + "\n" + userLine
}

func (m dashboardModel) renderMenuBar() string {
	divider := strings.Repeat("─", max(m.width-2, 40))
	return subtleStyle.Render(divider)
}

func (m dashboardModel) renderServerList() string {
	var b strings.Builder

	sectionTitle := titleStyle.Render("SERVERS")
	b.WriteString(sectionTitle)
	b.WriteString("\n")

	// Column headers
	header := fmt.Sprintf("  %-18s  %-16s  %-20s  %8s  %s",
		"NAME", "IP", "STATUS", "COST", "PROVIDER")
	b.WriteString(subtleStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", len(header))))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(subtleStyle.Render("Loading..."))
		return b.String()
	}

	if len(m.servers) == 0 {
		b.WriteString(subtleStyle.Render("No servers found. Press 'n' to create one."))
		return b.String()
	}

	for i, server := range m.servers {
		line := m.renderServerLine(server, i == m.cursor)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func (m dashboardModel) renderServerLine(server api.Server, selected bool) string {
	// Format each column
	name := truncate(server.Name, 18)
	ip := server.IPAddress
	if ip == "" {
		ip = "-"
	}
	ip = padRight(ip, 16)
	status := FormatStatus(server.Status, server.Pending, server.Power)
	cost := FormatMoney(server.CostSoFar)
	provider := server.Provider

	// Build line
	var prefix string
	if selected {
		prefix = selectedStyle.Render("> ")
		name = selectedStyle.Render(padRight(name, 18))
	} else {
		prefix = "  "
		name = padRight(name, 18)
	}

	// Pad status for alignment (account for ANSI codes)
	statusPadded := padRight(status, 20)

	line := fmt.Sprintf("%s%s  %s  %s  %s  %s",
		prefix, name, ip, statusPadded, padLeft(cost, 8), provider)

	return line
}

// Helper functions

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func padRight(s string, width int) string {
	// Account for ANSI escape codes
	visibleLen := lipgloss.Width(s)
	if visibleLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

func padLeft(s string, width int) string {
	visibleLen := lipgloss.Width(s)
	if visibleLen >= width {
		return s
	}
	return strings.Repeat(" ", width-visibleLen) + s
}
