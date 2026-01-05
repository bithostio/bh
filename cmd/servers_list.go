package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var serverListPage int

var serverListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all servers",
	RunE:    runServerList,
}

func init() {
	serverCmd.AddCommand(serverListCmd)
	serverListCmd.Flags().IntVar(&serverListPage, "page", 0, "Page number (optional)")
}

func runServerList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListServers(serverListPage)
	if err != nil {
		return err
	}

	displayServers(resp.Servers)
	ui.DisplayPagination(resp.Meta.Pagination)
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
		return ui.Yellow("pending")
	}

	switch server.Status {
	case "active":
		if server.Power {
			return ui.Green("active")
		}
		return ui.Yellow("powered off")
	case "failed":
		return ui.Red("failed")
	default:
		return server.Status
	}
}
