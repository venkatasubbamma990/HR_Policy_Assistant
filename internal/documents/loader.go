package documents

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"
)

// Loader reads and parses policy documents from a directory.
type Loader struct {
	dir string
	log *zap.Logger
}

// NewLoader creates a loader for the given documents directory.
func NewLoader(dir string, log *zap.Logger) *Loader {
	return &Loader{
		dir: dir,
		log: log.Named("documents"),
	}
}

// LoadAll reads and parses all markdown policy files.
func (l *Loader) LoadAll() ([]PolicyDocument, error) {
	l.log.Info("loading policy documents", zap.String("directory", l.dir))

	entries, err := os.ReadDir(l.dir)
	if err != nil {
		l.log.Error("failed to read documents directory", zap.String("directory", l.dir), zap.Error(err))
		return nil, fmt.Errorf("read dir %q: %w", l.dir, err)
	}

	var filenames []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		filenames = append(filenames, entry.Name())
	}
	sort.Strings(filenames)

	if len(filenames) == 0 {
		l.log.Warn("no markdown documents found", zap.String("directory", l.dir))
		return nil, fmt.Errorf("no markdown documents found in %q", l.dir)
	}

	l.log.Info("found policy files", zap.Int("count", len(filenames)), zap.Strings("files", filenames))

	docs := make([]PolicyDocument, 0, len(filenames))
	for _, name := range filenames {
		doc, err := l.LoadFile(name)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	l.log.Info("policy documents loaded successfully", zap.Int("count", len(docs)))
	return docs, nil
}

// LoadFile reads and parses a single policy file by name.
func (l *Loader) LoadFile(filename string) (PolicyDocument, error) {
	path := filepath.Join(l.dir, filename)

	content, err := os.ReadFile(path)
	if err != nil {
		l.log.Error("failed to read policy file", zap.String("path", path), zap.Error(err))
		return PolicyDocument{}, fmt.Errorf("read file %q: %w", path, err)
	}

	doc, err := ParsePolicyDocument(filename, path, string(content))
	if err != nil {
		l.log.Error("failed to parse policy file", zap.String("source", filename), zap.Error(err))
		return PolicyDocument{}, err
	}

	l.log.Debug("parsed policy file",
		zap.String("source", doc.Metadata.Source),
		zap.String("doc_id", doc.Metadata.DocID),
		zap.String("policy_type", doc.Metadata.PolicyType),
		zap.Int("sections", len(doc.Sections)),
	)

	return doc, nil
}

// ParsePolicyDocument parses raw markdown into a structured policy document.
func ParsePolicyDocument(sourceFilename, path, content string) (PolicyDocument, error) {
	meta, title, err := ParseMetadata(sourceFilename, content)
	if err != nil {
		return PolicyDocument{}, err
	}

	sections, err := ParseMarkdown(content)
	if err != nil {
		return PolicyDocument{}, fmt.Errorf("parse markdown %q: %w", sourceFilename, err)
	}

	return PolicyDocument{
		Metadata: meta,
		Title:    title,
		Path:     path,
		Sections: sections,
		Raw:      content,
	}, nil
}
