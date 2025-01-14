package main

import (
	"github.com/fwartner/dot/cmd"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dotfiles",
		Short: "Dotfiles Manager",
		Long:  "A tool to manage your dotfiles with ease.",
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewInstallCommand())
	rootCmd.AddCommand(cmd.NewSetupCommand())
	rootCmd.AddCommand(cmd.NewInitCommand())
	rootCmd.AddCommand(cmd.NewPullCommand())
	rootCmd.AddCommand(cmd.NewPushCommand())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Application failed: %v", err)
	}
}
