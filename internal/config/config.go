package config

import (
	"os"
	"path/filepath"
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

	return &Config{
		DocumentsDir:  documentsDir,
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		EmbedModel:    envOrDefault("HR_EMBED_MODEL", "text-embedding-3-small"),
		ChatModel:     envOrDefault("HR_CHAT_MODEL", "gpt-4o-mini"),
		IngestVerbose: os.Getenv("HR_INGEST_VERBOSE") == "1",
		LogLevel:      envOrDefault("LOG_LEVEL", "info"),
		LogFormat:     envOrDefault("LOG_FORMAT", "console"),
	}, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
