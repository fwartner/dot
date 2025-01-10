package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// InstallTools installs tools based on the current platform, optionally skipping specified tools
func InstallTools(skipTools []string) {
	// Convert the slice of tools to a set for faster lookup
	skipSet := make(map[string]struct{})
	for _, tool := range skipTools {
		skipSet[tool] = struct{}{}
	}

	logrus.Info("Installing necessary tools...")
	if runtime.GOOS == "linux" {
		distro := DetectLinuxDistro()
		logrus.Infof("Detected Linux distribution: %s", distro)
		switch distro {
		case "ubuntu", "debian":
			RunCommand("sudo", "apt", "update")
			for _, tool := range Config.Tools {
				if _, skipped := skipSet[tool]; skipped {
					logrus.Infof("Skipping tool: %s", tool)
					continue
				}
				if !IsToolInstalled(tool) {
					RunCommand("sudo", "apt", "install", "-y", tool)
				} else {
					logrus.Infof("Tool already installed: %s", tool)
				}
			}
		case "fedora":
			for _, tool := range Config.Tools {
				if _, skipped := skipSet[tool]; skipped {
					logrus.Infof("Skipping tool: %s", tool)
					continue
				}
				if !IsToolInstalled(tool) {
					RunCommand("sudo", "dnf", "install", "-y", tool)
				} else {
					logrus.Infof("Tool already installed: %s", tool)
				}
			}
		case "arch":
			EnsureYayInstalled()
			for _, tool := range Config.Tools {
				if _, skipped := skipSet[tool]; skipped {
					logrus.Infof("Skipping tool: %s", tool)
					continue
				}
				if !IsToolInstalled(tool) {
					RunCommand("yay", "-S", "--noconfirm", tool)
				} else {
					logrus.Infof("Tool already installed: %s", tool)
				}
			}
		default:
			logrus.Fatalf("Unsupported Linux distribution: %s", distro)
		}
	} else if runtime.GOOS == "darwin" {
		for _, tool := range Config.Tools {
			if _, skipped := skipSet[tool]; skipped {
				logrus.Infof("Skipping tool: %s", tool)
				continue
			}
			if !IsToolInstalled(tool) {
				RunCommand("brew", "install", tool)
			} else {
				logrus.Infof("Tool already installed: %s", tool)
			}
		}
	} else {
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

// EnsureYayInstalled ensures the AUR helper "yay" is installed on Arch Linux
func EnsureYayInstalled() {
	if IsToolInstalled("yay") {
		logrus.Info("Yay is already installed.")
		return
	}
	logrus.Info("Installing yay...")
	RunCommand("sudo", "pacman", "-Sy", "--noconfirm")
	RunCommand("git", "clone", "https://aur.archlinux.org/yay.git")
	RunCommand("sh", "-c", "cd yay && makepkg -si --noconfirm && cd .. && rm -rf yay")
}

// IsToolInstalled checks if a tool is installed
func IsToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
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
