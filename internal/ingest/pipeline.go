package ingest

import (
	"fmt"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/chunk"
	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/documents"
)

// Result holds the output of the full ingestion and chunking pipeline.
type Result struct {
	Documents []documents.PolicyDocument
	Chunks    []chunk.Chunk
}

// Pipeline loads, parses, and chunks HR policy documents for indexing.
type Pipeline struct {
	loader  *documents.Loader
	chunker *chunk.Chunker
	log     *zap.Logger
}

// NewPipeline creates an ingestion pipeline for the given documents directory.
func NewPipeline(cfg *config.Config, log *zap.Logger) *Pipeline {
	return &Pipeline{
		loader:  documents.NewLoader(cfg.DocumentsDir, log),
		chunker: chunk.NewChunker(cfg.ChunkConfig(), log),
		log:     log.Named("ingest"),
	}
}

// Run executes document load, parse, and chunking for indexing.
func (p *Pipeline) Run() (*Result, error) {
	p.log.Info("starting document ingestion pipeline")

	docs, err := p.loader.LoadAll()
	if err != nil {
		p.log.Error("document ingestion failed", zap.Error(err))
		return nil, fmt.Errorf("ingest: %w", err)
	}

	p.logResults(docs)

	p.log.Info("starting document chunking pipeline",
		zap.Int("document_count", len(docs)),
	)

	chunks := p.chunker.ChunkAll(docs)
	p.logChunkResults(chunks)

	return &Result{
		Documents: docs,
		Chunks:    chunks,
	}, nil
}

// LogVerboseDetails emits debug-level logs for each parsed section and chunk.
func (p *Pipeline) LogVerboseDetails(result *Result) {
	for _, doc := range result.Documents {
		p.log.Debug("policy document details",
			zap.String("title", doc.Title),
			zap.String("source", doc.Metadata.Source),
			zap.String("policy_type", doc.Metadata.PolicyType),
			zap.String("doc_id", doc.Metadata.DocID),
			zap.String("version", doc.Metadata.Version),
			zap.String("effective_date", doc.Metadata.EffectiveDate),
		)

		for _, section := range doc.Sections {
			p.log.Debug("parsed section",
				zap.String("source", doc.Metadata.Source),
				zap.Int("level", section.Level),
				zap.String("heading", section.Heading),
				zap.Int("tables", len(section.Tables)),
				zap.Int("lists", len(section.Lists)),
				zap.Int("paragraphs", len(section.Paragraphs)),
			)
		}
	}

	for _, c := range result.Chunks {
		p.log.Debug("indexed chunk",
			zap.String("chunk_id", c.ID),
			zap.String("source", c.Source),
			zap.String("section_path", c.SectionPath),
			zap.Int("token_count", c.TokenCount),
		)
	}
}

func (p *Pipeline) logResults(docs []documents.PolicyDocument) {
	p.log.Info("document ingestion completed", zap.Int("document_count", len(docs)))

	for _, doc := range docs {
		p.log.Info("indexed policy document",
			zap.String("source", doc.Metadata.Source),
			zap.String("policy_type", doc.Metadata.PolicyType),
			zap.String("doc_id", doc.Metadata.DocID),
			zap.String("version", doc.Metadata.Version),
			zap.String("effective_date", doc.Metadata.EffectiveDate),
			zap.Int("sections", len(doc.Sections)),
		)
	}
}

func (p *Pipeline) logChunkResults(chunks []chunk.Chunk) {
	p.log.Info("document chunking completed", zap.Int("chunk_count", len(chunks)))

	bySource := make(map[string]int)
	for _, c := range chunks {
		bySource[c.Source]++
	}

	for source, count := range bySource {
		p.log.Info("document chunks created",
			zap.String("source", source),
			zap.Int("chunk_count", count),
		)
	}
}
