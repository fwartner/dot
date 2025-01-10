package cmd

import (
	"dotfiles/utils"

	"github.com/spf13/cobra"
)

func NewUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the dotfiles manager tool",
		Run: func(cmd *cobra.Command, args []string) {
			utils.UpdateTool()
		},
	}
}
