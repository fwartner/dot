package utils

import (
	"os"
	"os/exec"
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

// PushDotfiles commits and pushes local changes to the dotfiles repository
func PushDotfiles() {
	logrus.Info("Pushing local changes to dotfiles repository")
	RunCommand("git", "-C", Config.DotfilesDir, "add", ".")
	RunCommand("git", "-C", Config.DotfilesDir, "commit", "-m", "Updated dotfiles")
	RunCommand("git", "-C", Config.DotfilesDir, "push")
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
