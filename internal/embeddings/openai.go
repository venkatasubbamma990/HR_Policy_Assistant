package embeddings

import (
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"

	"hrpolicyassistant/internal/config"
)

// NewOpenAIEmbedder creates a langchaingo embedder using the configured OpenAI model.
func NewOpenAIEmbedder(cfg *config.Config) (embeddings.Embedder, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	llm, err := openai.New(
		openai.WithToken(cfg.OpenAIAPIKey),
		openai.WithEmbeddingModel(cfg.EmbedModel),
		openai.WithEmbeddingDimensions(cfg.EmbedDimensions),
	)
	if err != nil {
		return nil, fmt.Errorf("create openai client: %w", err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, fmt.Errorf("create embedder: %w", err)
	}

	return embedder, nil
}
