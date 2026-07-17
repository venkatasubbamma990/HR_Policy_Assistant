package config

import (
	"bufio"
	"os"
	"strings"
)

// loadEnvFile reads a .env file and sets variables that are not already defined.
func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
