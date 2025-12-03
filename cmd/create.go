package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server (interactive wizard)",
	RunE:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	fmt.Println(ui.Bold("Bithost Server Creation Wizard"))
	fmt.Println(strings.Repeat("=", 50))

	// Step 1: Select Provider
	fmt.Println("\n" + ui.Cyan("Step 1: Select Provider"))
	providers, err := client.ListProviders()
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}
	provider, err := ui.PromptProvider(providers)
	if err != nil {
		return err
	}

	// Step 2: Select Region
	fmt.Println("\n" + ui.Cyan("Step 2: Select Region"))
	regions, err := client.ListRegions(provider.ID)
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}
	region, err := ui.PromptRegion(regions)
	if err != nil {
		return err
	}

	// Step 3: Select Size/Plan
	fmt.Println("\n" + ui.Cyan("Step 3: Select Size/Plan"))
	sizes, err := client.ListSizes(region.ID, provider.ID)
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}
	size, err := ui.PromptSize(sizes)
	if err != nil {
		return err
	}

	// Step 4: Select Operating System
	fmt.Println("\n" + ui.Cyan("Step 4: Select Operating System"))
	images, err := client.ListImages(provider.ID, "x86")
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}
	image, err := ui.PromptImage(images)
	if err != nil {
		return err
	}

	// Step 5: Select SSH Keys
	fmt.Println("\n" + ui.Cyan("Step 5: Select SSH Keys"))
	keys, err := client.ListSSHKeys()
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}
	selectedKeyIDs, err := ui.PromptSSHKeys(keys)
	if err != nil {
		return err
	}

	// Step 6: Enable Backups
	fmt.Println("\n" + ui.Cyan("Step 6: Backups"))
	backupsEnabled, err := ui.PromptBackups()
	if err != nil {
		return err
	}

	// Step 7: Server Name
	fmt.Println("\n" + ui.Cyan("Step 7: Server Name"))
	defaultName := fmt.Sprintf("server-%d", time.Now().Unix())
	serverName, err := ui.PromptServerName(defaultName)
	if err != nil {
		return err
	}

	// Summary
	summary := buildSummary(serverName, provider, region, size, image, backupsEnabled, len(selectedKeyIDs))
	confirmed, err := ui.ConfirmCreation(summary)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	// Create server
	stop := ui.ShowProgress("Creating server...")

	req := &api.CreateServerRequest{
		Name:           serverName,
		SizeID:         size.ID,
		RegionID:       region.ID,
		ImageID:        image.ID,
		ProviderID:     provider.ID,
		KeyIDs:         selectedKeyIDs,
		BackupsEnabled: backupsEnabled,
		Terms:          true, // Auto-accept terms
	}

	server, err := client.CreateServer(req)
	stop()

	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}

	fmt.Printf("\n%s Server created successfully!\n\n", ui.Green("✓"))
	fmt.Printf("  ID:     %d\n", server.ID)
	fmt.Printf("  Name:   %s\n", server.Name)
	fmt.Printf("  Status: %s\n", ui.Yellow("Pending"))
	fmt.Printf("\nYour server is being provisioned. Run 'bh servers' to check status.\n")

	return nil
}

func buildSummary(name string, provider *api.Provider, region *api.Region, size *api.Size, image *api.Image, backups bool, keyCount int) string {
	var sb strings.Builder

	sb.WriteString(ui.Bold("Server Configuration Summary:\n"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	sb.WriteString(fmt.Sprintf("  Name:         %s\n", name))
	sb.WriteString(fmt.Sprintf("  Provider:     %s\n", provider.Name))
	sb.WriteString(fmt.Sprintf("  Region:       %s\n", region.Name))
	sb.WriteString(fmt.Sprintf("  Size:         %s (%dMB RAM, %s)\n", size.Name, size.Memory, size.Processor))
	sb.WriteString(fmt.Sprintf("  OS:           %s\n", image.Name))
	sb.WriteString(fmt.Sprintf("  SSH Keys:     %d selected\n", keyCount))
	sb.WriteString(fmt.Sprintf("  Backups:      %s\n", formatBool(backups)))
	sb.WriteString(fmt.Sprintf("  Monthly Cost: %s", ui.FormatMoney(size.Price)))
	if backups {
		sb.WriteString(fmt.Sprintf(" + %s (backups)", ui.FormatMoney(size.Price*0.2)))
	}
	sb.WriteString("\n")

	return sb.String()
}

func formatBool(b bool) string {
	if b {
		return ui.Green("Enabled")
	}
	return "Disabled"
}
