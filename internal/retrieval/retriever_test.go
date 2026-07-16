package retrieval

import (
	"testing"

	"github.com/tmc/langchaingo/schema"
)

func TestToChunks(t *testing.T) {
	docs := []schema.Document{
		{
			PageContent: "Notice period for L3 is 60 days.",
			Score:       0.91,
			Metadata: map[string]any{
				"chunk_id":     "HR-POL-NOT-006-0001",
				"source":       "notice-policy.md",
				"policy_type":  "notice",
				"section_path": "Notice Period > L3",
				"document_id":  "HR-POL-NOT-006",
			},
		},
	}

	chunks := toChunks(docs)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0].Source != "notice-policy.md" {
		t.Fatalf("source = %q", chunks[0].Source)
	}
	if chunks[0].Score != 0.91 {
		t.Fatalf("score = %v", chunks[0].Score)
	}
}

func TestMetadataString(t *testing.T) {
	if got := metadataString(map[string]any{"policy_type": "notice"}, "policy_type"); got != "notice" {
		t.Fatalf("metadataString() = %q", got)
	}
}
