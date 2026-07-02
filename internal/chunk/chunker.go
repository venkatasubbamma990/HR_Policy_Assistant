package chunk

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/documents"
)

// Chunker splits policy documents into retrieval-sized chunks.
type Chunker struct {
	cfg Config
	log *zap.Logger
}

// NewChunker creates a chunker with the given configuration.
func NewChunker(cfg Config, log *zap.Logger) *Chunker {
	if log == nil {
		log = zap.NewNop()
	}
	return &Chunker{
		cfg: cfg,
		log: log.Named("chunk"),
	}
}

// ChunkAll splits all policy documents into chunks.
func (c *Chunker) ChunkAll(docs []documents.PolicyDocument) []Chunk {
	var chunks []Chunk
	for _, doc := range docs {
		chunks = append(chunks, c.ChunkDocument(doc)...)
	}
	return chunks
}

// ChunkDocument splits a single policy document by section headings.
func (c *Chunker) ChunkDocument(doc documents.PolicyDocument) []Chunk {
	var chunks []Chunk
	headingStack := make([]string, 0, 6)
	chunkCounter := 0

	for _, section := range doc.Sections {
		for len(headingStack) >= section.Level {
			headingStack = headingStack[:len(headingStack)-1]
		}
		headingStack = append(headingStack, section.Heading)

		sectionPath := strings.Join(headingStack, " > ")
		body := RenderSection(section)
		if body == "" {
			continue
		}

		prefix := c.sectionPrefix(doc, sectionPath)
		sectionChunks := c.splitSection(prefix, body)

		for _, content := range sectionChunks {
			tokenCount := EstimateTokens(content)
			chunk := Chunk{
				ID:             fmt.Sprintf("%s-%04d", doc.Metadata.DocID, chunkCounter),
				DocumentID:     doc.Metadata.DocID,
				Source:         doc.Metadata.Source,
				PolicyType:     doc.Metadata.PolicyType,
				SectionHeading: section.Heading,
				SectionPath:    sectionPath,
				Content:        content,
				TokenCount:     tokenCount,
				ChunkIndex:     chunkCounter,
			}
			chunks = append(chunks, chunk)
			chunkCounter++

			c.log.Debug("created chunk",
				zap.String("chunk_id", chunk.ID),
				zap.String("source", chunk.Source),
				zap.String("section_path", chunk.SectionPath),
				zap.Int("token_count", chunk.TokenCount),
			)
		}
	}

	c.log.Info("document chunked",
		zap.String("source", doc.Metadata.Source),
		zap.String("doc_id", doc.Metadata.DocID),
		zap.Int("chunk_count", len(chunks)),
	)

	return chunks
}

func (c *Chunker) sectionPrefix(doc documents.PolicyDocument, sectionPath string) string {
	return fmt.Sprintf(
		"Policy: %s (%s)\nDocument ID: %s | Version: %s | Effective: %s\nSection: %s\n\n",
		doc.Title,
		doc.Metadata.PolicyType,
		doc.Metadata.DocID,
		doc.Metadata.Version,
		doc.Metadata.EffectiveDate,
		sectionPath,
	)
}

func (c *Chunker) splitSection(prefix, body string) []string {
	full := strings.TrimSpace(prefix + body)
	if EstimateTokens(full) <= c.cfg.MaxTokens {
		return []string{full}
	}

	units := splitUnits(body)
	if len(units) == 0 {
		return []string{full}
	}

	var chunks []string
	var current strings.Builder
	current.WriteString(prefix)
	currentTokens := EstimateTokens(prefix)
	var overlap string

	flush := func() {
		text := strings.TrimSpace(current.String())
		if text != "" {
			chunks = append(chunks, text)
			overlap = OverlapText(text, c.cfg.OverlapTokens)
		}
		current.Reset()
		if overlap != "" {
			current.WriteString(prefix)
			current.WriteString(overlap)
			current.WriteString("\n\n")
			currentTokens = EstimateTokens(current.String())
		} else {
			current.WriteString(prefix)
			currentTokens = EstimateTokens(prefix)
		}
	}

	for _, unit := range units {
		unitTokens := EstimateTokens(unit)
		if unitTokens > c.cfg.MaxTokens {
			if currentTokens > EstimateTokens(prefix) {
				flush()
			}
			for _, part := range splitOversizedUnit(unit, c.cfg.MaxTokens-EstimateTokens(prefix)) {
				chunks = append(chunks, strings.TrimSpace(prefix+part))
			}
			overlap = ""
			current.Reset()
			current.WriteString(prefix)
			currentTokens = EstimateTokens(prefix)
			continue
		}

		if currentTokens+unitTokens > c.cfg.MaxTokens && currentTokens > EstimateTokens(prefix) {
			flush()
		}

		if current.Len() > len(prefix) {
			current.WriteString("\n\n")
		}
		current.WriteString(unit)
		currentTokens = EstimateTokens(current.String())
	}

	if currentTokens > EstimateTokens(prefix) {
		flush()
	} else if len(chunks) == 0 {
		chunks = append(chunks, strings.TrimSpace(current.String()))
	}

	return chunks
}

func splitUnits(body string) []string {
	parts := strings.Split(body, "\n\n")
	units := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			units = append(units, part)
		}
	}
	return units
}

func splitOversizedUnit(unit string, maxTokens int) []string {
	if maxTokens <= 0 {
		maxTokens = 200
	}

	sentences := splitSentences(unit)
	if len(sentences) == 0 {
		return []string{unit}
	}

	var parts []string
	var current strings.Builder
	currentTokens := 0

	flush := func() {
		if current.Len() > 0 {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			currentTokens = 0
		}
	}

	for _, sentence := range sentences {
		sentenceTokens := EstimateTokens(sentence)
		if sentenceTokens > maxTokens {
			flush()
			for _, wordPart := range splitByWords(sentence, maxTokens) {
				parts = append(parts, wordPart)
			}
			continue
		}

		if currentTokens+sentenceTokens > maxTokens && currentTokens > 0 {
			flush()
		}

		if current.Len() > 0 {
			current.WriteString(" ")
		}
		current.WriteString(sentence)
		currentTokens = EstimateTokens(current.String())
	}

	flush()
	return parts
}

func splitSentences(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	var sentences []string
	var current strings.Builder
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		current.WriteRune(runes[i])
		if runes[i] == '.' || runes[i] == '!' || runes[i] == '?' {
			if i+1 < len(runes) && runes[i+1] == ' ' {
				sentences = append(sentences, strings.TrimSpace(current.String()))
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(current.String()))
	}

	if len(sentences) == 0 {
		return []string{text}
	}
	return sentences
}

func splitByWords(text string, maxTokens int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var parts []string
	var current strings.Builder
	currentTokens := 0

	flush := func() {
		if current.Len() > 0 {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			currentTokens = 0
		}
	}

	for _, word := range words {
		wordTokens := EstimateTokens(word)
		if currentTokens+wordTokens > maxTokens && currentTokens > 0 {
			flush()
		}
		if current.Len() > 0 {
			current.WriteString(" ")
		}
		current.WriteString(word)
		currentTokens = EstimateTokens(current.String())
	}

	flush()
	return parts
}
