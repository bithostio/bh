package ui

import (
	"fmt"
	"time"

	"github.com/bithostio/bh/internal/api"
	"github.com/fatih/color"
)

var (
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)

// FormatMoney formats a float as currency
func FormatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// DisplayServers displays a table of servers
func DisplayServers(servers []api.Server) {
	if len(servers) == 0 {
		fmt.Println("No servers found.")
		return
	}

	// Print header
	fmt.Printf("\n%-8s %-20s %-15s %-20s %-10s %-10s\n", "ID", "NAME", "STATUS", "IP ADDRESS", "COST", "PROVIDER")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────")

	// Print rows
	for _, server := range servers {
		status := formatStatus(server)
		cost := FormatMoney(server.CostSoFar)
		ipAddr := server.IPAddress
		if ipAddr == "" {
			ipAddr = "-"
		}

		fmt.Printf("%-8d %-20s %-15s %-20s %-10s %-10d\n",
			server.ID,
			truncate(server.Name, 20),
			status,
			ipAddr,
			cost,
			server.ProviderID,
		)
	}
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatStatus(server api.Server) string {
	if server.Pending {
		return Yellow("Pending")
	}

	switch server.Status {
	case 1: // Active
		if server.Power == "on" {
			return Green("Active")
		}
		return Yellow("Powered Off")
	case 2: // Failed
		return Red("Failed")
	default:
		return fmt.Sprintf("Status %d", server.Status)
	}
}

// ShowProgress displays a spinner with a message
func ShowProgress(message string) func() {
	done := make(chan bool)

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r%s %s\n", Green("✓"), message)
				return
			default:
				fmt.Printf("\r%s %s", frames[i], message)
				i = (i + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		done <- true
		time.Sleep(200 * time.Millisecond)
	}
}
