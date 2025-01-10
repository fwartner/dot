package cmd

import (
	"dotfiles/utils"
	"strings"

	"github.com/spf13/cobra"
)

// NewInstallCommand creates the install command
func NewInstallCommand() *cobra.Command {
	var skipTools string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install necessary tools and dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			// Split the skipTools flag into a slice of tools
			skipped := strings.Split(skipTools, ",")
			utils.InstallTools(skipped)
		},
	}

	// Add a flag for skipping tools
	cmd.Flags().StringVar(&skipTools, "skip", "", "Comma-separated list of tools to skip during installation")

	return cmd
}
