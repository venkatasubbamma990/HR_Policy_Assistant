package embeddings

import (
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"

	"hrpolicyassistant/internal/config"
)

// NewOpenAIEmbedder creates a langchaingo embedder using an OpenAI-compatible endpoint.
// Works with OpenAI, LM Studio, and other OpenAI-compatible local servers.
func NewOpenAIEmbedder(cfg *config.Config) (embeddings.Embedder, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required (use lm-studio when using LM Studio)")
	}

	opts := []openai.Option{
		openai.WithToken(cfg.OpenAIAPIKey),
		openai.WithEmbeddingModel(cfg.EmbedModel),
	}
	if cfg.OpenAIBaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.OpenAIBaseURL))
	}
	if cfg.EmbedDimensions > 0 {
		opts = append(opts, openai.WithEmbeddingDimensions(cfg.EmbedDimensions))
	}

	llm, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create openai client: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("create embedder: %w", err)
	}

	return embedder, nil
}
