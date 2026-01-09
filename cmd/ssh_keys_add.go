package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/cli"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	keyLabel      string
	keyFile       string
	keyContent    string
	keyProviderID int
)

var sshKeysAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH key",
	Long: `Add a new SSH key for server access.

Examples:
  # Interactive mode with prompts
  bh ssh-keys add

  # From a file
  bh ssh-keys add --label "My Key" --file ~/.ssh/id_rsa.pub

  # Direct key content
  bh ssh-keys add --label "My Key" --key "ssh-rsa AAAA..."
`,
	RunE: runSSHKeysAdd,
}

func init() {
	sshKeysCmd.AddCommand(sshKeysAddCmd)
	sshKeysAddCmd.Flags().StringVarP(&keyLabel, "label", "l", "", "Label for the SSH key")
	sshKeysAddCmd.Flags().StringVarP(&keyFile, "file", "f", "", "Path to SSH public key file")
	sshKeysAddCmd.Flags().StringVarP(&keyContent, "key", "k", "", "SSH public key content")
	sshKeysAddCmd.Flags().IntVarP(&keyProviderID, "provider", "p", 1, "Provider ID (default: 1)")
}

func runSSHKeysAdd(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	label := keyLabel
	if label == "" {
		prompt := promptui.Prompt{
			Label: "SSH Key Label",
			Validate: func(input string) error {
				if len(input) < 1 {
					return fmt.Errorf("label cannot be empty")
				}
				return nil
			},
		}
		label, err = prompt.Run()
		if err != nil {
			return err
		}
	}

	var publicKey string
	if keyContent != "" {
		publicKey = keyContent
	} else if keyFile != "" {
		data, err := os.ReadFile(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read key file: %w", err)
		}
		publicKey = strings.TrimSpace(string(data))
	} else {
		prompt := promptui.Prompt{
			Label:   "Path to SSH public key file",
			Default: "~/.ssh/id_rsa.pub",
		}
		path, err := prompt.Run()
		if err != nil {
			return err
		}

		if strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = strings.Replace(path, "~", home, 1)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read key file: %w", err)
		}
		publicKey = strings.TrimSpace(string(data))
	}

	if publicKey == "" {
		return fmt.Errorf("SSH key content is empty")
	}

	stop, cancel := cli.ShowProgress("Adding SSH key...")

	req := &api.CreateSSHKeyRequest{
		Label:      label,
		Key:        publicKey,
		ProviderID: keyProviderID,
	}

	key, err := client.CreateSSHKey(req)

	if err != nil {
		cancel()
		return err
	}

	stop()

	fmt.Printf("\n%s SSH key added successfully!\n\n", cli.Green("✓"))
	fmt.Printf("  ID:    %d\n", key.ID)
	fmt.Printf("  Label: %s\n", key.Label)

	return nil
}
