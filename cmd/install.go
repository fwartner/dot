package cmd

import (
	"github.com/fwartner/dot/utils"
	"github.com/spf13/cobra"
)

// NewInstallCommand creates the install command
func NewInstallCommand() *cobra.Command {
	var skipTools string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install necessary tools and dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			utils.InstallTools()
		},
	}

	// Add a flag for skipping tools
	cmd.Flags().StringVar(&skipTools, "skip", "", "Comma-separated list of tools to skip during installation")

	return cmd
}
