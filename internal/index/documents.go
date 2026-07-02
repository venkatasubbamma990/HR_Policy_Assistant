package index

import (
	"hrpolicyassistant/internal/chunk"

	"github.com/tmc/langchaingo/schema"
)

const indexBatchSize = 32

// ChunksToDocuments converts application chunks into langchaingo documents.
func ChunksToDocuments(chunks []chunk.Chunk) []schema.Document {
	docs := make([]schema.Document, 0, len(chunks))
	for _, c := range chunks {
		docs = append(docs, schema.Document{
			PageContent: c.Content,
			Metadata: map[string]any{
				"chunk_id":        c.ID,
				"document_id":     c.DocumentID,
				"source":          c.Source,
				"policy_type":     c.PolicyType,
				"section_heading": c.SectionHeading,
				"section":         c.SectionPath,
				"section_path":    c.SectionPath,
				"token_count":     c.TokenCount,
				"chunk_index":     c.ChunkIndex,
			},
		})
	}
	return docs
}
