package llm

import (
	"fmt"
	"strings"

	"hrpolicyassistant/internal/config"
)

// WrapAPIReachabilityError adds setup hints when an OpenAI-compatible endpoint is unreachable.
func WrapAPIReachabilityError(cfg *config.Config, err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()
	if !strings.Contains(msg, "failed to reach API server") &&
		!strings.Contains(msg, "connection refused") &&
		!strings.Contains(msg, "connect: network is unreachable") {
		return err
	}

	baseURL := cfg.OpenAIBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	hint := fmt.Sprintf("cannot reach LLM/embeddings API at %s", baseURL)
	hint += "; " + config.WSLSetupHint(baseURL)

	return fmt.Errorf("%s: %w", hint, err)
}
