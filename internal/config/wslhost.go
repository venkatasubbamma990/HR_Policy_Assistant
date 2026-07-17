package config

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const defaultWSLProxyPort = 12345

// applyWSLOpenAIBaseURL rewrites localhost LM Studio URLs when running inside WSL.
// Requires scripts/setup-wsl-llm.ps1 (Windows portproxy 12345 -> 1234) or mirrored networking.
func applyWSLOpenAIBaseURL() {
	if !RunningInWSL() || os.Getenv("HR_WSL_LLM_DISABLE") == "1" {
		return
	}

	raw := os.Getenv("OPENAI_BASE_URL")
	if raw == "" {
		raw = "http://127.0.0.1:1234/v1"
	}
	if !isLocalhostURL(raw) {
		return
	}

	gateway := wslDefaultGateway()
	if gateway == "" {
		return
	}

	proxyPort := envOrDefault("HR_WSL_LLM_PORT", strconv.Itoa(defaultWSLProxyPort))
	resolved, err := rewriteURLHost(raw, gateway, proxyPort)
	if err != nil {
		return
	}

	_ = os.Setenv("OPENAI_BASE_URL", resolved)
}

func isLocalhostURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.Contains(raw, "127.0.0.1") || strings.Contains(raw, "localhost")
	}

	host := strings.ToLower(parsed.Hostname())
	return host == "127.0.0.1" || host == "localhost"
}

func wslDefaultGateway() string {
	if host := strings.TrimSpace(os.Getenv("WSL_HOST_IP")); host != "" {
		return host
	}

	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[1] != "00000000" {
			continue
		}

		gatewayHex := fields[2]
		if len(gatewayHex) != 8 {
			continue
		}

		value, err := strconv.ParseUint(gatewayHex, 16, 32)
		if err != nil || value == 0 {
			continue
		}

		ip := make(net.IP, 4)
		binary.LittleEndian.PutUint32(ip, uint32(value))
		return ip.String()
	}

	return ""
}

func rewriteURLHost(raw, host, port string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	parsed.Host = net.JoinHostPort(host, port)
	return parsed.String(), nil
}

// WSLSetupHint returns setup guidance when LM Studio is unreachable from WSL.
func WSLSetupHint(baseURL string) string {
	if !RunningInWSL() {
		return fmt.Sprintf("ensure LM Studio is running (Developer → Start Server) and OPENAI_BASE_URL (%s) is correct", baseURL)
	}

	return "run scripts/setup-wsl-llm.ps1 once in PowerShell as Administrator (forwards port 12345 to LM Studio), or enable WSL mirrored networking in %USERPROFILE%\\.wslconfig ([wsl2] networkingMode=mirrored) and run wsl --shutdown"
}
