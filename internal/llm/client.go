package llm

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"hrpolicyassistant/internal/config"
)

// NewChatModel creates an OpenAI-compatible chat model client.
func NewChatModel(cfg *config.Config) (llms.Model, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	opts := []openai.Option{
		openai.WithToken(cfg.OpenAIAPIKey),
		openai.WithModel(cfg.ChatModel),
	}
	if cfg.OpenAIBaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.OpenAIBaseURL))
	}

	model, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create chat model: %w", err)
	}

	return model, nil
}
