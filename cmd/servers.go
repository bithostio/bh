package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var serversCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"ls", "list"},
	Short:   "List all your servers",
	RunE:    runServers,
}

func init() {
	rootCmd.AddCommand(serversCmd)
}

func runServers(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	servers, err := client.ListServers()
	if err != nil {
		return err
	}

	displayServers(servers)

	return nil
}

func displayServers(servers []api.Server) {
	if len(servers) == 0 {
		fmt.Println("No servers found.")
		return
	}

	rows := [][]string{
		{"ID", "Name", "Status", "IP Address", "Cost", "Provider"},
	}

	for _, server := range servers {
		ipAddr := server.IPAddress
		if ipAddr == "" {
			ipAddr = "-"
		}

		rows = append(rows, []string{
			strconv.Itoa(server.ID),
			ui.Truncate(server.Name, 20),
			formatServerStatus(server),
			ipAddr,
			ui.FormatMoney(server.CostSoFar),
			strconv.Itoa(server.ProviderID),
		})
	}

	ui.PrintTable(rows)
}

func formatServerStatus(server api.Server) string {
	if server.Pending {
		return ui.Yellow("Pending")
	}

	switch server.Status {
	case 1: // Active
		if server.Power == "on" {
			return ui.Green("Active")
		}
		return ui.Yellow("Powered Off")
	case 2: // Failed
		return ui.Red("Failed")
	default:
		return fmt.Sprintf("Status %d", server.Status)
	}
}
