package chunk

import (
	"strings"
	"testing"

	"hrpolicyassistant/internal/documents"
)

func TestEstimateTokens(t *testing.T) {
	got := EstimateTokens("hello world")
	if got < 1 {
		t.Fatalf("expected positive token count, got %d", got)
	}
}

func TestOverlapText(t *testing.T) {
	text := "Privilege Leave can be taken during notice period only with manager approval and HR clearance."
	overlap := OverlapText(text, 10)
	if overlap == "" {
		t.Fatal("expected overlap text")
	}
	if EstimateTokens(overlap) > 15 {
		t.Fatalf("overlap too large: %q (%d tokens)", overlap, EstimateTokens(overlap))
	}
}

func TestChunkDocumentBySection(t *testing.T) {
	doc := sampleDocument()
	chunker := NewChunker(DefaultConfig(), nil)
	chunks := chunker.ChunkDocument(doc)

	if len(chunks) == 0 {
		t.Fatal("expected chunks")
	}

	seen := make(map[string]bool)
	for _, c := range chunks {
		if c.DocumentID != "HR-POL-LVE-002" {
			t.Fatalf("unexpected document id %q", c.DocumentID)
		}
		if c.SectionPath == "" {
			t.Fatal("expected section path")
		}
		if c.TokenCount <= 0 {
			t.Fatal("expected positive token count")
		}
		if c.TokenCount > DefaultConfig().MaxTokens+200 {
			t.Fatalf("chunk exceeds max tokens: %d", c.TokenCount)
		}
		if seen[c.ID] {
			t.Fatalf("duplicate chunk id %q", c.ID)
		}
		seen[c.ID] = true
	}
}

func TestSplitLargeSectionWithOverlap(t *testing.T) {
	cfg := Config{MinTokens: 50, MaxTokens: 120, OverlapTokens: 20}
	chunker := NewChunker(cfg, nil)

	body := strings.Repeat("Employees may apply for Privilege Leave through HRMS. ", 20)
	chunks := chunker.splitSection("Prefix:\n\n", body)

	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	for i, c := range chunks {
		tokens := EstimateTokens(c)
		if tokens > cfg.MaxTokens+50 {
			t.Fatalf("chunk %d too large: %d tokens", i, tokens)
		}
	}
}

func TestRenderSectionIncludesTableAndList(t *testing.T) {
	section := documents.Section{
		Level:   3,
		Heading: "Privilege Leave",
		Paragraphs: []string{
			"Privilege Leave (PL) is available to confirmed employees.",
		},
		Tables: []documents.Table{
			{Headers: []string{"Leave Type", "Days/Year"}, Rows: [][]string{{"PL", "18"}}},
		},
		Lists: []documents.List{
			{Items: []string{"Apply via HRMS"}},
		},
	}

	text := RenderSection(section)
	for _, want := range []string{"Privilege Leave", "Days/Year", "Apply via HRMS"} {
		if !strings.Contains(text, want) {
			t.Fatalf("rendered section missing %q: %q", want, text)
		}
	}
}

func TestChunkAllRealDocuments(t *testing.T) {
	dir := "../../documents"
	loader := documents.NewLoader(dir, nil)
	docs, err := loader.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}

	chunker := NewChunker(DefaultConfig(), nil)
	chunks := chunker.ChunkAll(docs)
	if len(chunks) == 0 {
		t.Fatal("expected chunks from real documents")
	}

	cfg := DefaultConfig()
	for _, c := range chunks {
		if c.TokenCount > cfg.MaxTokens+250 {
			t.Fatalf("chunk %s exceeds max tokens: %d", c.ID, c.TokenCount)
		}
		if !strings.Contains(c.Content, c.DocumentID) {
			t.Fatalf("chunk %s missing document id in content", c.ID)
		}
	}
}

func sampleDocument() documents.PolicyDocument {
	return documents.PolicyDocument{
		Metadata: documents.Metadata{
			Source:        "leave-policy.md",
			PolicyType:    "leave",
			DocID:         "HR-POL-LVE-002",
			Version:       "4.1",
			EffectiveDate: "April 1, 2025",
		},
		Title: "Leave Policy",
		Sections: []documents.Section{
			{Level: 2, Heading: "Leave Types", Paragraphs: []string{"Overview of leave types."}},
			{
				Level:      3,
				Heading:    "Privilege Leave",
				Paragraphs: []string{"PL requires manager approval."},
				Lists: []documents.List{
					{Items: []string{"Apply via HRMS", "Minimum 3 days notice"}},
				},
			},
		},
	}
}
