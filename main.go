package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"dotfiles/cmd"
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
	rootCmd.AddCommand(cmd.NewPullCommand())
	rootCmd.AddCommand(cmd.NewPushCommand())
	rootCmd.AddCommand(cmd.NewUpdateCommand())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Application failed: %v", err)
	}
}
