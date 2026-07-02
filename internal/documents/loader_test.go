package documents

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestParseMetadata(t *testing.T) {
	content := `# Leave Policy

**Company:** NovaTech Solutions Pvt. Ltd.  
**Document ID:** HR-POL-LVE-002  
**Version:** 4.1  
**Effective Date:** April 1, 2025  
**Applies To:** All permanent employees

---
`

	meta, title, err := ParseMetadata("leave-policy.md", content)
	if err != nil {
		t.Fatalf("ParseMetadata() error = %v", err)
	}

	if title != "Leave Policy" {
		t.Fatalf("title = %q, want %q", title, "Leave Policy")
	}
	if meta.Source != "leave-policy.md" {
		t.Fatalf("source = %q", meta.Source)
	}
	if meta.PolicyType != "leave" {
		t.Fatalf("policy_type = %q, want leave", meta.PolicyType)
	}
	if meta.DocID != "HR-POL-LVE-002" {
		t.Fatalf("doc_id = %q", meta.DocID)
	}
	if meta.Version != "4.1" {
		t.Fatalf("version = %q", meta.Version)
	}
	if meta.EffectiveDate != "April 1, 2025" {
		t.Fatalf("effective_date = %q", meta.EffectiveDate)
	}
}

func TestPolicyTypeFromSource(t *testing.T) {
	cases := map[string]string{
		"wfh-policy.md":     "wfh",
		"salary-policy.md":  "salary",
		"notice-policy.md":  "notice",
		"entry-policy.md":   "entry",
		"exit-policy.md":    "exit",
		"leave-policy.md":   "leave",
	}

	for source, want := range cases {
		got := policyTypeFromSource(source)
		if got != want {
			t.Fatalf("policyTypeFromSource(%q) = %q, want %q", source, got, want)
		}
	}
}

func TestLoadAllRealDocuments(t *testing.T) {
	dir := filepath.Join("..", "..", "documents")
	loader := NewLoader(dir, zap.NewNop())
	docs, err := loader.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}

	if len(docs) != 6 {
		t.Fatalf("loaded %d docs, want 6", len(docs))
	}

	for _, doc := range docs {
		if doc.Metadata.DocID == "" {
			t.Fatalf("missing doc_id for %s", doc.Metadata.Source)
		}
		if doc.Metadata.PolicyType == "" {
			t.Fatalf("missing policy_type for %s", doc.Metadata.Source)
		}
		if len(doc.Sections) == 0 {
			t.Fatalf("expected parsed sections for %s", doc.Metadata.Source)
		}
	}
}

func TestParseMarkdownTablesAndLists(t *testing.T) {
	content := `# Test Policy

**Document ID:** HR-POL-TST-001  
**Version:** 1.0  
**Effective Date:** Jan 1, 2025  

---

## Benefits

| Type | Days |
|------|------|
| PL | 18 |
| CL | 12 |

- Apply via HRMS
- Manager approval required

1. Submit request
2. Wait for approval
`

	sections, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if len(sections) == 0 {
		t.Fatal("expected sections")
	}

	section := sections[0]
	if section.Heading != "Benefits" {
		t.Fatalf("heading = %q", section.Heading)
	}
	if len(section.Tables) != 1 {
		t.Fatalf("tables = %d, want 1", len(section.Tables))
	}
	if len(section.Tables[0].Headers) != 2 {
		t.Fatalf("table headers = %d", len(section.Tables[0].Headers))
	}
	if len(section.Lists) != 2 {
		t.Fatalf("lists = %d, want 2", len(section.Lists))
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
