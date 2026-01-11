package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/cli"
	"github.com/spf13/cobra"
)

var (
	serverListPage int
	showAllServers bool
)

var serverListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List active servers (use --all for all servers)",
	RunE:    runServerList,
}

func init() {
	serverCmd.AddCommand(serverListCmd)
	serverListCmd.Flags().IntVar(&serverListPage, "page", 0, "Page number (optional)")
	serverListCmd.Flags().BoolVarP(&showAllServers, "all", "a", false, "Show all servers including failed")
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

	servers := resp.Servers
	// Filter to active servers by default
	if !showAllServers {
		servers = filterActiveServers(servers)
	}

	displayServers(servers, showAllServers)
	cli.DisplayPagination(resp.Meta.Pagination)
	return nil
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

func displayServers(servers []api.Server, showAll bool) {
	if len(servers) == 0 {
		if showAll {
			fmt.Println("No servers found.")
		} else {
			fmt.Println("No active servers found. Use --all to see all servers.")
		}
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
			cli.Truncate(server.Name, 20),
			formatServerStatus(server),
			ipAddr,
			cli.FormatMoney(server.CostSoFar),
			strconv.Itoa(server.ProviderID),
		})
	}

	cli.PrintTable(rows)
}

func formatServerStatus(server api.Server) string {
	if server.Pending {
		return cli.Yellow("pending")
	}

	switch server.Status {
	case "active":
		if server.Power {
			return cli.Green("active")
		}
		return cli.Yellow("powered off")
	case "failed":
		return cli.Red("failed")
	default:
		return server.Status
	}
}
