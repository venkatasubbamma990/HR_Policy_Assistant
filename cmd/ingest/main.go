package main

import (
	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/ingest"
	"hrpolicyassistant/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("config: " + err.Error())
	}

	log, err := logger.New(logger.Options{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	if err != nil {
		panic("logger: " + err.Error())
	}
	defer func() { _ = log.Sync() }()

	log = log.Named("ingest")
	log.Info("starting document ingestion",
		zap.String("documents_dir", cfg.DocumentsDir),
		zap.Bool("verbose", cfg.IngestVerbose),
		zap.Int("chunk_min_tokens", cfg.ChunkMinTokens),
		zap.Int("chunk_max_tokens", cfg.ChunkMaxTokens),
		zap.Int("chunk_overlap_tokens", cfg.ChunkOverlapTokens),
	)

	pipeline := ingest.NewPipeline(cfg, log)
	result, err := pipeline.Run()
	if err != nil {
		log.Fatal("document ingestion failed", zap.Error(err))
	}

	if cfg.IngestVerbose {
		pipeline.LogVerboseDetails(result)
	}

	log.Info("document ingestion finished",
		zap.Int("document_count", len(result.Documents)),
		zap.Int("chunk_count", len(result.Chunks)),
	)
}
