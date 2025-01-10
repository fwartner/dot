package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Configuration structure for config.yml
type Configuration struct {
	DotfilesRepo string   `yaml:"dotfiles_repo"`
	DotfilesDir  string   `yaml:"dotfiles_dir"`
	Tools        []string `yaml:"tools"`
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

// InstallTools installs tools based on the current platform, skipping specified tools
func InstallTools(skipTools []string) {
	skipSet := make(map[string]struct{})
	for _, tool := range skipTools {
		skipSet[tool] = struct{}{}
	}

	logrus.Info("Installing necessary tools...")
	for _, tool := range Config.Tools {
		if _, skipped := skipSet[tool]; skipped {
			logrus.Infof("Skipping tool: %s", tool)
			continue
		}
		if !IsToolInstalled(tool) {
			logrus.Infof("Installing tool: %s", tool)
			installTool(tool)
		} else {
			logrus.Infof("Tool already installed: %s", tool)
		}
	}
}

// IsToolInstalled checks if a tool is installed
func IsToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// installTool installs a tool based on the current platform
func installTool(tool string) {
	switch runtime.GOOS {
	case "linux":
		distro := DetectLinuxDistro()
		switch distro {
		case "ubuntu", "debian":
			RunCommand("sudo", "apt", "update")
			RunCommand("sudo", "apt", "install", "-y", tool)
		case "fedora":
			RunCommand("sudo", "dnf", "install", "-y", tool)
		case "arch":
			if !IsToolInstalled("yay") {
				logrus.Info("Installing yay...")
				RunCommand("sudo", "pacman", "-Sy", "--noconfirm")
				RunCommand("git", "clone", "https://aur.archlinux.org/yay.git")
				RunCommand("sh", "-c", "cd yay && makepkg -si --noconfirm && cd .. && rm -rf yay")
			}
			RunCommand("yay", "-S", "--noconfirm", tool)
		default:
			logrus.Fatalf("Unsupported Linux distribution: %s", distro)
		}
	case "darwin":
		RunCommand("brew", "install", tool)
	default:
		logrus.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}
}

// DetectLinuxDistro detects the current Linux distribution
func DetectLinuxDistro() string {
	cmd := exec.Command("sh", "-c", "cat /etc/os-release | grep ^ID=")
	output, err := cmd.Output()
	if err != nil {
		logrus.Fatalf("Failed to detect Linux distribution: %v", err)
	}
	return strings.TrimSpace(strings.Split(string(output), "=")[1])
}

// CloneDotfiles clones the dotfiles repository
func CloneDotfiles() {
	if _, err := os.Stat(Config.DotfilesDir); err == nil {
		logrus.Infof("Dotfiles directory already exists: %s", Config.DotfilesDir)
		return
	}

	logrus.Infof("Cloning dotfiles repository: %s", Config.DotfilesRepo)
	RunCommand("git", "clone", Config.DotfilesRepo, Config.DotfilesDir)
}

// StowDotfiles uses GNU Stow to symlink dotfiles
func StowDotfiles() {
	logrus.Infof("Stowing dotfiles from: %s", Config.DotfilesDir)
	RunCommand("stow", "-d", Config.DotfilesDir, "-t", os.ExpandEnv("$HOME"), ".")
}

// PullDotfiles pulls the latest changes from the dotfiles repository
func PullDotfiles() {
	logrus.Info("Pulling latest changes from dotfiles repository")
	RunCommand("git", "-C", Config.DotfilesDir, "pull")
}

// UpdateTool updates the tool by pulling the latest version and rebuilding
func UpdateTool() {
	logrus.Info("Updating the dotfiles manager tool...")
	RunCommand("git", "pull", "origin", "main")
	RunCommand("go", "build", "-o", "/usr/local/bin/dotfiles")
	logrus.Info("Dotfiles manager tool updated successfully.")
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

// PushDotfiles commits and pushes local changes to the dotfiles repository
func PushDotfiles() {
	dotfilesDir := Config.DotfilesDir
	if dotfilesDir == "" {
		dotfilesDir = os.ExpandEnv("$HOME/dotfiles")
	}

	// Check if there are changes to commit
	output, err := RunGitCommand(dotfilesDir, "status", "--porcelain")
	if err != nil {
		logrus.Fatalf("Failed to check git status: %v", err)
	}

	if strings.TrimSpace(output) == "" {
		logrus.Info("No changes to commit.")
		return
	}

	// Stage changes
	logrus.Info("Staging changes...")
	if _, err := RunGitCommand(dotfilesDir, "add", "."); err != nil {
		logrus.Fatalf("Failed to stage changes: %v", err)
	}

	// Generate commit message based on changes
	commitMessage := GenerateCommitMessage(dotfilesDir)
	logrus.Infof("Committing changes with message: %s", commitMessage)
	if _, err := RunGitCommand(dotfilesDir, "commit", "-m", commitMessage); err != nil {
		logrus.Fatalf("Failed to commit changes: %v", err)
	}

	// Push changes
	logrus.Info("Pushing changes to the repository...")
	if _, err := RunGitCommand(dotfilesDir, "push"); err != nil {
		logrus.Fatalf("Failed to push changes: %v", err)
	}

	logrus.Info("Changes pushed successfully.")
}

// GenerateCommitMessage creates a commit message based on file changes
func GenerateCommitMessage(repoDir string) string {
	output, err := RunGitCommand(repoDir, "diff", "--cached", "--name-status")
	if err != nil {
		logrus.Fatalf("Failed to generate commit message: %v", err)
	}

	var added, modified, deleted []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		file := parts[1]

		switch status {
		case "A":
			added = append(added, file)
		case "M":
			modified = append(modified, file)
		case "D":
			deleted = append(deleted, file)
		}
	}

	// Generate a summary message
	var messageParts []string
	if len(added) > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Added: %s", strings.Join(added, ", ")))
	}
	if len(modified) > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Modified: %s", strings.Join(modified, ", ")))
	}
	if len(deleted) > 0 {
		messageParts = append(messageParts, fmt.Sprintf("Deleted: %s", strings.Join(deleted, ", ")))
	}

	if len(messageParts) == 0 {
		return "Updated dotfiles"
	}

	return strings.Join(messageParts, "; ")
}

// InitDotfilesRepo initializes a new dotfiles repository
func InitDotfilesRepo(remote string) error {
	// Use the configured directory or fallback to the default
	dotfilesDir := Config.DotfilesDir
	if dotfilesDir == "" {
		dotfilesDir = os.ExpandEnv("$HOME/dotfiles")
	}

	// Create the dotfiles directory if it doesn't exist
	if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
		logrus.Infof("Creating dotfiles directory: %s", dotfilesDir)
		if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
			return fmt.Errorf("failed to create dotfiles directory: %w", err)
		}
	}

	// Initialize Git repository
	logrus.Infof("Initializing Git repository in: %s", dotfilesDir)
	cmd := exec.Command("git", "init")
	cmd.Dir = dotfilesDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	// Add a remote repository if provided
	if remote != "" {
		logrus.Infof("Adding remote origin: %s", remote)
		cmd := exec.Command("git", "remote", "add", "origin", remote)
		cmd.Dir = dotfilesDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add remote repository: %w", err)
		}

		// Fetch the remote's default branch
		logrus.Info("Fetching remote default branch...")
		cmd = exec.Command("git", "fetch", "--all")
		cmd.Dir = dotfilesDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to fetch remote branches: %w", err)
		}

		// Determine the default branch of the remote
		output, err := RunGitCommand(dotfilesDir, "symbolic-ref", "refs/remotes/origin/HEAD")
		if err != nil {
			logrus.Warn("Unable to determine remote default branch. Defaulting to 'main'.")
		} else {
			defaultBranch := strings.TrimSpace(strings.Replace(output, "refs/remotes/origin/", "", 1))
			logrus.Infof("Setting default branch to: %s", defaultBranch)

			// Set the default branch locally
			cmd = exec.Command("git", "branch", "-M", defaultBranch)
			cmd.Dir = dotfilesDir
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to set default branch to %s: %w", defaultBranch, err)
			}
		}
	}

	// Create a basic .gitignore file
	gitignorePath := filepath.Join(dotfilesDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		logrus.Infof("Creating .gitignore file in: %s", dotfilesDir)
		content := `
# Ignore stow-related files
*/.stow-local-ignore
*/.stow-global-ignore
`
		if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore file: %w", err)
		}
	}

	logrus.Infof("Dotfiles repository initialized at: %s", dotfilesDir)
	return nil
}

// RunGitCommand runs a git command in the specified repository directory
func RunGitCommand(repoDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git command failed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
