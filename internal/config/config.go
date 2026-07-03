package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"hrpolicyassistant/internal/chunk"
)

const defaultDocumentsDir = "documents"

// Config holds application settings for the HR Policy RAG assistant.
type Config struct {
	DocumentsDir       string
	OpenAIAPIKey       string
	EmbedModel         string
	EmbedDimensions    int
	ChatModel          string
	DatabaseURL        string
	PGVectorCollection string
	IndexForce         bool
	IngestVerbose      bool
	LogLevel           string
	LogFormat          string
	ChunkMinTokens     int
	ChunkMaxTokens     int
	ChunkOverlapTokens int
	HTTPPort           int
}

// ChunkConfig returns chunking settings from configuration.
func (c *Config) ChunkConfig() chunk.Config {
	return chunk.Config{
		MinTokens:     c.ChunkMinTokens,
		MaxTokens:     c.ChunkMaxTokens,
		OverlapTokens: c.ChunkOverlapTokens,
	}
}

// IndexingEnabled reports whether vector indexing is configured.
func (c *Config) IndexingEnabled() bool {
	return c.DatabaseURL != "" && c.OpenAIAPIKey != ""
}

// HTTPAddr returns the query API listen address.
func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
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

	cfg := &Config{
		DocumentsDir:       documentsDir,
		OpenAIAPIKey:       os.Getenv("OPENAI_API_KEY"),
		EmbedModel:         envOrDefault("HR_EMBED_MODEL", "text-embedding-3-small"),
		EmbedDimensions:    envIntOrDefault("HR_EMBED_DIMENSIONS", 1536),
		ChatModel:          envOrDefault("HR_CHAT_MODEL", "gpt-4o-mini"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		PGVectorCollection: envOrDefault("PGVECTOR_COLLECTION", "hr_policies"),
		IndexForce:         os.Getenv("HR_INDEX_FORCE") == "1",
		IngestVerbose:      os.Getenv("HR_INGEST_VERBOSE") == "1",
		LogLevel:           envOrDefault("LOG_LEVEL", "info"),
		LogFormat:          envOrDefault("LOG_FORMAT", "console"),
		ChunkMinTokens:     envIntOrDefault("HR_CHUNK_MIN_TOKENS", defaults.MinTokens),
		ChunkMaxTokens:     envIntOrDefault("HR_CHUNK_MAX_TOKENS", defaults.MaxTokens),
		ChunkOverlapTokens: envIntOrDefault("HR_CHUNK_OVERLAP_TOKENS", defaults.OverlapTokens),
		HTTPPort:           envIntOrDefault("HR_HTTP_PORT", 8080),
	}

	if cfg.EmbedDimensions <= 0 {
		return nil, fmt.Errorf("HR_EMBED_DIMENSIONS must be positive")
	}

	return cfg, nil
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
