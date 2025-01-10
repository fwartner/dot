package cmd

import (
	"github.com/fwartner/dot/utils"

	"github.com/spf13/cobra"
)

func NewPullCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull the latest changes from the dotfiles repository",
		Run: func(cmd *cobra.Command, args []string) {
			utils.PullDotfiles()
		},
	}
}
