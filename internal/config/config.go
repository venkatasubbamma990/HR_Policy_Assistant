package config

import (
	"os"
	"path/filepath"
	"strconv"

	"hrpolicyassistant/internal/chunk"
)

const defaultDocumentsDir = "documents"

// Config holds application settings for the HR Policy RAG assistant.
type Config struct {
	DocumentsDir  string
	OpenAIAPIKey  string
	EmbedModel    string
	ChatModel     string
	IngestVerbose bool
	LogLevel      string
	LogFormat     string
	ChunkMinTokens     int
	ChunkMaxTokens     int
	ChunkOverlapTokens int
}

// ChunkConfig returns chunking settings from configuration.
func (c *Config) ChunkConfig() chunk.Config {
	return chunk.Config{
		MinTokens:     c.ChunkMinTokens,
		MaxTokens:     c.ChunkMaxTokens,
		OverlapTokens: c.ChunkOverlapTokens,
	}
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	documentsDir := os.Getenv("HR_DOCUMENTS_DIR")
	if documentsDir == "" {
		documentsDir = filepath.Join(wd, defaultDocumentsDir)
	}

	defaults := chunk.DefaultConfig()

	return &Config{
		DocumentsDir:       documentsDir,
		OpenAIAPIKey:       os.Getenv("OPENAI_API_KEY"),
		EmbedModel:         envOrDefault("HR_EMBED_MODEL", "text-embedding-3-small"),
		ChatModel:          envOrDefault("HR_CHAT_MODEL", "gpt-4o-mini"),
		IngestVerbose:      os.Getenv("HR_INGEST_VERBOSE") == "1",
		LogLevel:           envOrDefault("LOG_LEVEL", "info"),
		LogFormat:          envOrDefault("LOG_FORMAT", "console"),
		ChunkMinTokens:     envIntOrDefault("HR_CHUNK_MIN_TOKENS", defaults.MinTokens),
		ChunkMaxTokens:     envIntOrDefault("HR_CHUNK_MAX_TOKENS", defaults.MaxTokens),
		ChunkOverlapTokens: envIntOrDefault("HR_CHUNK_OVERLAP_TOKENS", defaults.OverlapTokens),
	}, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
