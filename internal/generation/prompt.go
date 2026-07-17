package generation

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"hrpolicyassistant/internal/retrieval"
)

var citationPattern = regexp.MustCompile(`\[(\d+)\]`)

const systemPrompt = `You are an HR Policy Assistant for NovaTech Solutions.
Answer employee questions using ONLY the numbered policy excerpts provided.
Rules:
- Base the answer strictly on the excerpts; do not invent policy details.
- Cite supporting excerpts inline using [1], [2], etc.
- If the excerpts do not contain enough information, say so clearly and suggest contacting HR.
- Be concise, accurate, and professional.`

// Citation references a policy excerpt used in the answer.
type Citation struct {
	Index       int     `json:"index"`
	Source      string  `json:"source"`
	DocumentID  string  `json:"document_id"`
	SectionPath string  `json:"section_path"`
	ChunkID     string  `json:"chunk_id"`
	Score       float32 `json:"score,omitempty"`
}

// Result is a generated answer with source citations.
type Result struct {
	Text      string     `json:"text"`
	ChatModel string     `json:"chat_model"`
	Citations []Citation `json:"citations"`
}

// BuildUserPrompt formats retrieved chunks and the question for the LLM.
func BuildUserPrompt(question string, chunks []retrieval.Chunk) string {
	var b strings.Builder
	b.WriteString("Policy excerpts:\n\n")

	for i, chunk := range chunks {
		fmt.Fprintf(
			&b,
			"[%d] Source: %s | Doc ID: %s | Section: %s\n%s\n\n",
			i+1,
			chunk.Source,
			chunk.DocumentID,
			chunk.SectionPath,
			chunk.Content,
		)
	}

	fmt.Fprintf(&b, "Question: %s", question)
	return b.String()
}

// ExtractCitations maps inline [n] references to retrieved chunk metadata.
func ExtractCitations(answer string, chunks []retrieval.Chunk) []Citation {
	if len(chunks) == 0 {
		return nil
	}

	indexes := referencedIndexes(answer)
	citations := make([]Citation, 0, len(indexes))
	for _, index := range indexes {
		if index < 1 || index > len(chunks) {
			continue
		}
		chunk := chunks[index-1]
		citations = append(citations, Citation{
			Index:       index,
			Source:      chunk.Source,
			DocumentID:  chunk.DocumentID,
			SectionPath: chunk.SectionPath,
			ChunkID:     chunk.ChunkID,
			Score:       chunk.Score,
		})
	}

	return citations
}

func referencedIndexes(answer string) []int {
	matches := citationPattern.FindAllStringSubmatch(answer, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[int]struct{})
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		index, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		seen[index] = struct{}{}
	}

	indexes := make([]int, 0, len(seen))
	for index := range seen {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	return indexes
}

const noContextAnswer = "I could not find relevant HR policy sections for your question. Please contact HR for assistance."

// NoContextResult returns a standard response when retrieval finds no chunks.
func NoContextResult(chatModel string) *Result {
	return &Result{
		Text:      noContextAnswer,
		ChatModel: chatModel,
		Citations: nil,
	}
}
