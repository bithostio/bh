package cmd

import (
	"github.com/spf13/cobra"
)

var sshKeysCmd = &cobra.Command{
	Use:     "ssh-keys",
	Aliases: []string{"keys"},
	Short:   "Manage SSH keys",
	Long:    `List and add SSH keys for server access.`,
}

func init() {
	rootCmd.AddCommand(sshKeysCmd)
}
