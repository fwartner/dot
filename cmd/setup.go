package cmd

import (
	"github.com/fwartner/dot/utils"
	"github.com/spf13/cobra"
)

// NewSetupCommand creates the setup command
func NewSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Clone and stow dotfiles",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CloneDotfiles()
			utils.StowDotfiles()
		},
	}
}
