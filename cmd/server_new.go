package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	interactive    bool
	serverName     string
	providerID     int
	regionID       int
	sizeID         int
	imageID        int
	sshKeyIDs      []int
	backupsEnabled bool
)

var serverNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new server",
	Long: `Create a new server either interactively or with command-line flags.

Interactive mode (wizard):
  bh server new --interactive

Programmatic mode (with flags):
  bh server new --name myserver --provider 1 --region 2 --size 5 --image 10 --keys 1,2

List available resources:
  bh providers
  bh regions --provider 1
  bh sizes --provider 1 --region 2
  bh images --provider 1
  bh ssh-keys`,
	RunE: runServerNew,
}

func init() {
	serverCmd.AddCommand(serverNewCmd)

	serverNewCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive wizard mode")
	serverNewCmd.Flags().StringVarP(&serverName, "name", "n", "", "Server name")
	serverNewCmd.Flags().IntVarP(&providerID, "provider", "p", 0, "Provider ID")
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

	// Interactive mode
	if interactive {
		return runInteractiveServerCreation(client)
	}

	// Programmatic mode - validate required flags
	if serverName == "" {
		return fmt.Errorf("--name is required (or use --interactive)")
	}
	if providerID == 0 {
		return fmt.Errorf("--provider is required (or use --interactive)")
	}
	if regionID == 0 {
		return fmt.Errorf("--region is required (or use --interactive)")
	}
	if sizeID == 0 {
		return fmt.Errorf("--size is required (or use --interactive)")
	}
	if imageID == 0 {
		return fmt.Errorf("--image is required (or use --interactive)")
	}
	if len(sshKeyIDs) == 0 {
		return fmt.Errorf("--keys is required (or use --interactive)")
	}

	// Create server with provided flags
	return createServer(client, serverName, providerID, regionID, sizeID, imageID, sshKeyIDs, backupsEnabled)
}

func createServer(client *api.Client, name string, providerID, regionID, sizeID, imageID int, keyIDs []int, backups bool) error {
	stop := ui.ShowProgress("Creating server...")

	req := &api.CreateServerRequest{
		Name:           name,
		SizeID:         sizeID,
		RegionID:       regionID,
		ImageID:        imageID,
		ProviderID:     providerID,
		KeyIDs:         keyIDs,
		BackupsEnabled: backups,
		Terms:          true,
	}

	server, err := client.CreateServer(req)
	stop()

	if err != nil {
		return err
	}

	fmt.Printf("\n%s Server created successfully!\n\n", ui.Green("✓"))
	fmt.Printf("  ID:     %d\n", server.ID)
	fmt.Printf("  Name:   %s\n", server.Name)
	fmt.Printf("  Status: %s\n", ui.Yellow("Pending"))
	fmt.Printf("\nYour server is being provisioned. Run 'bh server list' to check status.\n")

	return nil
}

func runInteractiveServerCreation(client *api.Client) error {
	fmt.Println(ui.Header("Bithost Server Creation Wizard"))
	fmt.Println(ui.Divider(50))

	// Step 1: Select Provider
	fmt.Println("\n" + ui.StepHeader(1, "Select Provider"))
	providers, err := client.ListProviders()
	if err != nil {
		return err
	}
	provider, err := promptProvider(providers)
	if err != nil {
		return err
	}

	// Step 2: Select Region
	fmt.Println("\n" + ui.StepHeader(2, "Select Region"))
	regions, err := client.ListRegions(provider.ID)
	if err != nil {
		return err
	}
	region, err := promptRegion(regions)
	if err != nil {
		return err
	}

	// Step 3: Select Size/Plan
	fmt.Println("\n" + ui.StepHeader(3, "Select Size/Plan"))
	sizes, err := client.ListSizes(region.ID, provider.ID)
	if err != nil {
		return err
	}
	size, err := promptSize(sizes)
	if err != nil {
		return err
	}

	// Step 4: Select Operating System
	fmt.Println("\n" + ui.StepHeader(4, "Select Operating System"))
	images, err := client.ListImages(provider.ID, "x86")
	if err != nil {
		return err
	}
	image, err := promptImage(images)
	if err != nil {
		return err
	}

	// Step 5: Select SSH Keys
	fmt.Println("\n" + ui.StepHeader(5, "Select SSH Keys"))
	keys, err := client.ListSSHKeys()
	if err != nil {
		return err
	}
	selectedKeyIDs, err := promptSSHKeys(keys)
	if err != nil {
		return err
	}

	// Step 6: Enable Backups
	fmt.Println("\n" + ui.StepHeader(6, "Backups"))
	backups, err := promptBackups()
	if err != nil {
		return err
	}

	// Step 7: Server Name
	fmt.Println("\n" + ui.StepHeader(7, "Server Name"))
	defaultName := fmt.Sprintf("server-%d", time.Now().Unix())
	name, err := promptServerName(defaultName)
	if err != nil {
		return err
	}

	// Summary
	summary := buildSummary(name, provider, region, size, image, backups, len(selectedKeyIDs))
	confirmed, err := confirmCreation(summary)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	return createServer(client, name, provider.ID, region.ID, size.ID, image.ID, selectedKeyIDs, backups)
}

func promptProvider(providers []api.Provider) (*api.Provider, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "\U0001F4E1 {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select Provider",
		Items:     providers,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &providers[i], nil
}

func promptRegion(regions []api.Region) (*api.Region, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Origin }})",
		Inactive: "  {{ .Name }} ({{ .Origin }})",
		Selected: "\U0001F30D {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select Region",
		Items:     regions,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &regions[i], nil
}

func promptSize(sizes []api.Size) (*api.Size, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} - ${{ printf \"%.2f\" .Price }}/mo ({{ .Memory }}MB RAM, {{ .Processor }})",
		Inactive: "  {{ .Name }} - ${{ printf \"%.2f\" .Price }}/mo ({{ .Memory }}MB RAM, {{ .Processor }})",
		Selected: "\U0001F4BE {{ .Name | green }} - ${{ printf \"%.2f\" .Price }}/mo",
	}

	prompt := promptui.Select{
		Label:     "Select Size/Plan",
		Items:     sizes,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &sizes[i], nil
}

func promptImage(images []api.Image) (*api.Image, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "\U0001F4BF {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select Operating System",
		Items:     images,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &images[i], nil
}

func promptSSHKeys(keys []api.SSHKey) ([]int, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no SSH keys found. Add keys at https://dashboard.bithost.io/keys")
	}

	fmt.Println("\nAvailable SSH Keys:")
	for i, key := range keys {
		fmt.Printf("  [%d] %s\n", i+1, key.Label)
	}

	prompt := promptui.Prompt{
		Label:   "Select SSH keys (comma-separated numbers, or 'all')",
		Default: "all",
	}

	result, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	result = strings.TrimSpace(result)
	if result == "all" {
		ids := make([]int, len(keys))
		for i, key := range keys {
			ids[i] = key.ID
		}
		return ids, nil
	}

	var selectedIDs []int
	for _, s := range strings.Split(result, ",") {
		s = strings.TrimSpace(s)
		idx, err := strconv.Atoi(s)
		if err != nil || idx < 1 || idx > len(keys) {
			return nil, fmt.Errorf("invalid selection: %s", s)
		}
		selectedIDs = append(selectedIDs, keys[idx-1].ID)
	}

	return selectedIDs, nil
}

func promptBackups() (bool, error) {
	prompt := promptui.Select{
		Label: "Enable automatic backups? (+20% cost)",
		Items: []string{"No", "Yes"},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func promptServerName(defaultName string) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Server name",
		Default: defaultName,
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("server name cannot be empty")
			}
			if len(input) > 255 {
				return fmt.Errorf("server name too long")
			}
			return nil
		},
	}

	return prompt.Run()
}

func confirmCreation(summary string) (bool, error) {
	fmt.Println("\n" + summary)

	prompt := promptui.Prompt{
		Label:     "Create this server",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func buildSummary(name string, provider *api.Provider, region *api.Region, size *api.Size, image *api.Image, backups bool, keyCount int) string {
	var sb strings.Builder

	sb.WriteString(ui.Header("Server Configuration Summary:") + "\n")
	sb.WriteString(ui.Divider(50) + "\n")
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "Name:", name))
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "Provider:", provider.Name))
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "Region:", region.Name))
	sb.WriteString(fmt.Sprintf("  %-15s %s (%sMB RAM, %s)\n", "Size:", size.Name, size.Memory, size.Processor))
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "OS:", image.Name))
	sb.WriteString(fmt.Sprintf("  %-15s %d selected\n", "SSH Keys:", keyCount))
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "Backups:", ui.FormatEnabled(backups)))

	cost := ui.FormatMoney(size.Price)
	if backups {
		cost += fmt.Sprintf(" + %s (backups)", ui.FormatMoney(size.Price*0.2))
	}
	sb.WriteString(fmt.Sprintf("  %-15s %s\n", "Monthly Cost:", cost))

	return sb.String()
}
