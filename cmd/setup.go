package cmd

import (
	"github.com/fwartner/dot/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Clone and stow dotfiles",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CloneDotfiles()
			utils.StowDotfiles()
			logrus.Info("Dotfiles setup is complete!")
		},
	}
}
