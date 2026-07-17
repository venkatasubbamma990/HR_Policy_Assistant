package config

import (
	"os"
	"strings"
)

// RunningInWSL reports whether the process is running inside WSL (not Docker).
func RunningInWSL() bool {
	if runningInContainer() {
		return false
	}

	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	version := strings.ToLower(string(data))
	return strings.Contains(version, "microsoft") || strings.Contains(version, "wsl")
}

func runningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, "docker") || strings.Contains(content, "containerd")
}
