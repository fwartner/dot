package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Configuration structure for config.yml
type Configuration struct {
	DotfilesRepo string   `yaml:"dotfiles_repo"`
	DotfilesDir  string   `yaml:"dotfiles_dir"`
	Tools        []string `yaml:"tools"`
	Files        []string `yaml:"files"` // New field for specific files or directories
}

// Config holds the loaded configuration
var Config Configuration

// LoadConfig loads the configuration from the config.yml file
func LoadConfig() {
	configPaths := []string{
		"./config.yml",
		os.ExpandEnv("$HOME/.config/dotfiles/config.yml"),
		os.ExpandEnv("$HOME/.dotfiles-config.yml"),
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		logrus.Fatalf("No configuration file found in: %v", configPaths)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		logrus.Fatalf("Failed to read config file %s: %v", configFile, err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		logrus.Fatalf("Failed to parse config file %s: %v", configFile, err)
	}

	logrus.Infof("Configuration loaded from %s", configFile)

	// Set default values if not provided
	if Config.DotfilesDir == "" {
		Config.DotfilesDir = os.ExpandEnv("$HOME/dotfiles")
	}
	if Config.DotfilesRepo == "" {
		logrus.Fatalf("Dotfiles repository URL is required in the configuration")
	}
}

// InitDotfilesRepo initializes a new dotfiles repository
func InitDotfilesRepo(remote string) error {
	dotfilesDir := os.ExpandEnv(Config.DotfilesDir)
	if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
		logrus.Infof("Creating dotfiles directory: %s", dotfilesDir)
		if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
			return fmt.Errorf("failed to create dotfiles directory: %w", err)
		}
	}

	logrus.Infof("Initializing Git repository in: %s", dotfilesDir)
	cmd := exec.Command("git", "init")
	cmd.Dir = dotfilesDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	if remote != "" {
		logrus.Infof("Adding remote origin: %s", remote)
		cmd := exec.Command("git", "remote", "add", "origin", remote)
		cmd.Dir = dotfilesDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add remote repository: %w", err)
		}
	}
	return nil
}

// CloneDotfiles clones the dotfiles repository
func CloneDotfiles() {
	dotfilesDir := os.ExpandEnv(Config.DotfilesDir)
	if _, err := os.Stat(dotfilesDir); err == nil {
		logrus.Infof("Dotfiles directory already exists: %s", dotfilesDir)
		return
	}

	logrus.Infof("Cloning dotfiles repository: %s", Config.DotfilesRepo)
	RunCommand("git", "clone", Config.DotfilesRepo, dotfilesDir)
}

// StowDotfiles uses GNU Stow to symlink specific dotfiles
func StowDotfiles() {
	if len(Config.Files) == 0 {
		logrus.Info("No specific dotfiles configured for syncing. Stowing everything in the repository.")
		runStow(".")
		return
	}

	logrus.Info("Syncing specific dotfiles...")
	for _, file := range Config.Files {
		if err := runStow(file); err != nil {
			logrus.Fatalf("Failed to stow %s: %v", file, err)
		}
	}
	logrus.Info("Dotfiles synced successfully.")
}

// PullDotfiles pulls the latest changes from the dotfiles repository
func PullDotfiles() {
	dotfilesDir := os.ExpandEnv(Config.DotfilesDir)
	logrus.Info("Pulling latest changes from the repository...")
	RunCommand("git", "-C", dotfilesDir, "pull")
}

// runStow executes the stow command for a specific file or directory
func runStow(target string) error {
	dotfilesDir := os.ExpandEnv(Config.DotfilesDir)
	if _, err := os.Stat(filepath.Join(dotfilesDir, target)); os.IsNotExist(err) {
		return fmt.Errorf("target %s does not exist in %s", target, dotfilesDir)
	}

	cmd := exec.Command("stow", "-d", dotfilesDir, "-t", os.ExpandEnv("$HOME"), target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommand runs a system command and handles errors
func RunCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Command failed: %s %s: %v", command, strings.Join(args, " "), err)
	}
}
