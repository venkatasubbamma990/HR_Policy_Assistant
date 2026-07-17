package config

import (
	"os"
	"strings"
)

// RunningInWSL reports whether the process is running inside WSL.
func RunningInWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	version := strings.ToLower(string(data))
	return strings.Contains(version, "microsoft") || strings.Contains(version, "wsl")
}
