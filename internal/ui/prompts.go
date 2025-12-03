package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bithostio/bh/internal/api"
	"github.com/manifoldco/promptui"
)

// PromptProvider prompts user to select a provider
func PromptProvider(providers []api.Provider) (*api.Provider, error) {
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

// PromptRegion prompts user to select a region
func PromptRegion(regions []api.Region) (*api.Region, error) {
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

// PromptSize prompts user to select a size/plan
func PromptSize(sizes []api.Size) (*api.Size, error) {
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

// PromptImage prompts user to select an OS image
func PromptImage(images []api.Image) (*api.Image, error) {
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

// PromptSSHKeys prompts user to select SSH keys
func PromptSSHKeys(keys []api.SSHKey) ([]int, error) {
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

// PromptBackups prompts user to enable backups
func PromptBackups() (bool, error) {
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

// PromptServerName prompts user for server name
func PromptServerName(defaultName string) (string, error) {
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

// ConfirmCreation prompts user to confirm server creation
func ConfirmCreation(summary string) (bool, error) {
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
