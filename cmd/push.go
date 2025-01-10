package cmd

import (
	"dotfiles/utils"

	"github.com/spf13/cobra"
)

func NewPushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Push local changes to the dotfiles repository",
		Run: func(cmd *cobra.Command, args []string) {
			utils.PushDotfiles()
		},
	}
}
