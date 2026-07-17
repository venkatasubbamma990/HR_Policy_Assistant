package generation

import (
	"strings"
	"testing"

	"hrpolicyassistant/internal/retrieval"
)

func TestBuildUserPromptIncludesNumberedChunks(t *testing.T) {
	chunks := []retrieval.Chunk{
		{
			Content:     "Notice period for L3 is 60 days.",
			Source:      "notice-policy.md",
			DocumentID:  "HR-POL-NOT-006",
			SectionPath: "Notice Period > L3",
		},
	}

	prompt := BuildUserPrompt("What is the notice period for L3?", chunks)
	for _, want := range []string{"[1]", "notice-policy.md", "HR-POL-NOT-006", "Question: What is the notice period for L3?"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestExtractCitations(t *testing.T) {
	chunks := []retrieval.Chunk{
		{Source: "notice-policy.md", DocumentID: "HR-POL-NOT-006", SectionPath: "A", ChunkID: "c1", Score: 0.9},
		{Source: "leave-policy.md", DocumentID: "HR-POL-LVE-002", SectionPath: "B", ChunkID: "c2", Score: 0.8},
	}

	answer := "The notice period is 60 days [1]. PL rules may also apply [2][1]."
	citations := ExtractCitations(answer, chunks)

	if len(citations) != 2 {
		t.Fatalf("expected 2 citations, got %d", len(citations))
	}
	if citations[0].Source != "notice-policy.md" || citations[1].Source != "leave-policy.md" {
		t.Fatalf("unexpected citations: %#v", citations)
	}
}

func TestNoContextResult(t *testing.T) {
	result := NoContextResult("gpt-4o-mini")
	if result.Text == "" || result.ChatModel != "gpt-4o-mini" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
