package cmd

import (
	"dotfiles/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewInitCommand() *cobra.Command {
	var remote string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new dotfiles repository",
		Long:  "Sets up a new Git repository for managing dotfiles and prepares a directory structure.",
		Run: func(cmd *cobra.Command, args []string) {
			err := utils.InitDotfilesRepo(remote)
			if err != nil {
				fmt.Printf("Error initializing dotfiles repository: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Dotfiles repository initialized successfully!")
		},
	}

	// Add flag for specifying a remote repository URL
	cmd.Flags().StringVar(&remote, "remote", "", "Remote repository URL to add as origin")

	return cmd
}
