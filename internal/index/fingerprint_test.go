package index

import (
	"testing"

	"hrpolicyassistant/internal/documents"
)

func TestFingerprintStableOrdering(t *testing.T) {
	docsA := []documents.PolicyDocument{
		{Metadata: documents.Metadata{Source: "b.md"}, Raw: "beta"},
		{Metadata: documents.Metadata{Source: "a.md"}, Raw: "alpha"},
	}
	docsB := []documents.PolicyDocument{
		{Metadata: documents.Metadata{Source: "a.md"}, Raw: "alpha"},
		{Metadata: documents.Metadata{Source: "b.md"}, Raw: "beta"},
	}

	if Fingerprint(docsA) != Fingerprint(docsB) {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestFingerprintChangesWhenContentChanges(t *testing.T) {
	base := []documents.PolicyDocument{
		{Metadata: documents.Metadata{Source: "leave-policy.md"}, Raw: "version 1"},
	}
	updated := []documents.PolicyDocument{
		{Metadata: documents.Metadata{Source: "leave-policy.md"}, Raw: "version 2"},
	}

	if Fingerprint(base) == Fingerprint(updated) {
		t.Fatal("fingerprint should change when content changes")
	}
}
