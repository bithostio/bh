package cmd

import (
	"fmt"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/cli"
	"github.com/bithostio/bh/internal/config"
	"github.com/spf13/cobra"
)

var (
	serverName     string
	provider       string
	regionID       int
	sizeID         int
	imageID        int
	sshKeyIDs      []int
	backupsEnabled bool
)

var serverNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new server",
	Long: `Create a new server with command-line flags.

For interactive server creation, run 'bh' without arguments to launch the TUI.

Usage:
  bh servers new --name myserver --provider digital_ocean --region 2 --size 5 --image 10 --keys 1,2

List available resources:
  bh providers
  bh regions --provider digital_ocean
  bh sizes --provider digital_ocean --region 2
  bh images --provider digital_ocean
  bh ssh-keys`,
	RunE: runServerNew,
}

func init() {
	serverCmd.AddCommand(serverNewCmd)

	serverNewCmd.Flags().StringVarP(&serverName, "name", "n", "", "Server name")
	serverNewCmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider slug")
	serverNewCmd.Flags().IntVarP(&regionID, "region", "r", 0, "Region ID")
	serverNewCmd.Flags().IntVarP(&sizeID, "size", "s", 0, "Size/plan ID")
	serverNewCmd.Flags().IntVarP(&imageID, "image", "m", 0, "Image/OS ID")
	serverNewCmd.Flags().IntSliceVarP(&sshKeyIDs, "keys", "k", []int{}, "SSH key IDs (comma-separated)")
	serverNewCmd.Flags().BoolVarP(&backupsEnabled, "backups", "b", false, "Enable automatic backups")
}

func runServerNew(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	if serverName == "" {
		return fmt.Errorf("--name is required (or run 'bh' for interactive mode)")
	}
	if provider == "" {
		return fmt.Errorf("--provider is required (or run 'bh' for interactive mode)")
	}
	if regionID == 0 {
		return fmt.Errorf("--region is required (or run 'bh' for interactive mode)")
	}
	if sizeID == 0 {
		return fmt.Errorf("--size is required (or run 'bh' for interactive mode)")
	}
	if imageID == 0 {
		return fmt.Errorf("--image is required (or run 'bh' for interactive mode)")
	}
	if len(sshKeyIDs) == 0 {
		return fmt.Errorf("--keys is required (or run 'bh' for interactive mode)")
	}

	return createServer(client, serverName, provider, regionID, sizeID, imageID, sshKeyIDs, backupsEnabled)
}

func createServer(client *api.Client, name string, provider string, regionID, sizeID, imageID int, keyIDs []int, backups bool) error {
	stop, cancel := cli.ShowProgress("Creating server...")

	req := &api.CreateServerRequest{
		Name:           name,
		SizeID:         sizeID,
		RegionID:       regionID,
		ImageID:        imageID,
		Provider:       provider,
		KeyIDs:         keyIDs,
		BackupsEnabled: backups,
		Terms:          true,
	}

	server, err := client.CreateServer(req)

	if err != nil {
		cancel()
		return err
	}

	stop()

	fmt.Printf("\n%s Server created successfully!\n\n", cli.Green("✓"))
	fmt.Printf("  ID:     %d\n", server.ID)
	fmt.Printf("  Name:   %s\n", server.Name)
	fmt.Printf("  Status: %s\n", cli.Yellow("Pending"))
	fmt.Printf("\nYour server is being provisioned. Run 'bh servers list' to check status.\n")

	return nil
}
