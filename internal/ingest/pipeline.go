package ingest

import (
	"fmt"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/documents"
)

// Pipeline loads and parses HR policy documents for indexing.
type Pipeline struct {
	loader *documents.Loader
	log    *zap.Logger
}

// NewPipeline creates an ingestion pipeline for the given documents directory.
func NewPipeline(documentsDir string, log *zap.Logger) *Pipeline {
	return &Pipeline{
		loader: documents.NewLoader(documentsDir, log),
		log:    log.Named("ingest"),
	}
}

// Run executes the document ingestion (indexing) pipeline.
func (p *Pipeline) Run() ([]documents.PolicyDocument, error) {
	p.log.Info("starting document ingestion pipeline")

	docs, err := p.loader.LoadAll()
	if err != nil {
		p.log.Error("document ingestion failed", zap.Error(err))
		return nil, fmt.Errorf("ingest: %w", err)
	}

	p.logResults(docs)
	return docs, nil
}

// LogVerboseDetails emits debug-level logs for each parsed section.
func (p *Pipeline) LogVerboseDetails(docs []documents.PolicyDocument) {
	for _, doc := range docs {
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
