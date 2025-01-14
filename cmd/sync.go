package cmd

import (
	"github.com/fwartner/dot/utils"

	"github.com/spf13/cobra"
)

// NewSyncCommand creates the sync command
func NewSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Pull and stow the latest dotfiles",
		Long:  "Fetch the latest changes from the dotfiles repository and sync them locally using GNU Stow.",
		Run: func(cmd *cobra.Command, args []string) {
			utils.PullDotfiles()
			utils.StowDotfiles()
		},
	}
}
