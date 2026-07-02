package chunk

import (
	"fmt"
	"strings"

	"hrpolicyassistant/internal/documents"
)

// Chunk is a retrieval-ready slice of policy text with metadata.
type Chunk struct {
	ID             string `json:"id"`
	DocumentID     string `json:"document_id"`
	Source         string `json:"source"`
	PolicyType     string `json:"policy_type"`
	SectionHeading string `json:"section_heading"`
	SectionPath    string `json:"section_path"`
	Content        string `json:"content"`
	TokenCount     int    `json:"token_count"`
	ChunkIndex     int    `json:"chunk_index"`
}

// RenderSection converts a parsed section into plain text for chunking.
func RenderSection(section documents.Section) string {
	var b strings.Builder

	if section.Heading != "" && section.Heading != "Document Body" {
		fmt.Fprintf(&b, "%s\n\n", section.Heading)
	}

	for _, paragraph := range section.Paragraphs {
		fmt.Fprintf(&b, "%s\n\n", paragraph)
	}

	for _, table := range section.Tables {
		b.WriteString(renderTable(table))
		b.WriteString("\n")
	}

	for _, list := range section.Lists {
		b.WriteString(renderList(list))
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

func renderTable(table documents.Table) string {
	if len(table.Headers) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(strings.Join(table.Headers, " | "))
	b.WriteString("\n")

	for _, row := range table.Rows {
		b.WriteString(strings.Join(row, " | "))
		b.WriteString("\n")
	}

	return b.String()
}

func renderList(list documents.List) string {
	var b strings.Builder
	for i, item := range list.Items {
		if list.Ordered {
			fmt.Fprintf(&b, "%d. %s\n", i+1, item)
		} else {
			fmt.Fprintf(&b, "- %s\n", item)
		}
	}
	return b.String()
}
