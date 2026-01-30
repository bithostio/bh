package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/bithostio/bh/internal/api"
)

type serverDetailModel struct {
	server api.Server
}

func newServerDetailModel(server api.Server) serverDetailModel {
	return serverDetailModel{
		server: server,
	}
}

func (m serverDetailModel) View() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render(fmt.Sprintf("SERVER: %s", m.server.Name))
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	// Server details
	b.WriteString(m.renderField("Status", FormatStatus(m.server.Status, m.server.Pending, m.server.Power)))
	b.WriteString(m.renderField("IP Address", formatIP(m.server.IPAddress)))
	b.WriteString(m.renderField("Private IP", formatIP(m.server.PrivateIPAddress)))
	b.WriteString(m.renderField("IPv6", formatIP(m.server.IPAddressV6)))
	b.WriteString(m.renderField("Provider", m.server.Provider))
	b.WriteString(m.renderField("Cost So Far", FormatMoney(m.server.CostSoFar)))
	b.WriteString(m.renderField("Backups", FormatEnabled(m.server.BackupsEnabled)))

	if m.server.Message != "" {
		b.WriteString("\n")
		b.WriteString(m.renderField("Message", m.server.Message))
	}

	// SSH command
	if m.server.IPAddress != "" {
		b.WriteString("\n")
		b.WriteString(m.renderField("SSH", m.sshCommand()))
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(DetailHelp())

	return b.String()
}

func (m serverDetailModel) renderField(label, value string) string {
	return fmt.Sprintf("%s %s\n",
		labelStyle.Render(label+":"),
		valueStyle.Render(value))
}

func formatIP(ip string) string {
	if ip == "" {
		return subtleStyle.Render("-")
	}
	return ip
}

func (m serverDetailModel) sshCommand() string {
	return fmt.Sprintf("ssh root@%s", m.server.IPAddress)
}

func (m serverDetailModel) CopySSHCommand() error {
	return clipboard.WriteAll(m.sshCommand())
}
