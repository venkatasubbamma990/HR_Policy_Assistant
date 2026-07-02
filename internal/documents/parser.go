package documents

import (
	"regexp"
	"strings"
)

var headingPattern = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
var orderedListPattern = regexp.MustCompile(`^\d+\.\s+(.+)$`)

// PolicyDocument is a fully parsed HR policy file with structured content.
type PolicyDocument struct {
	Metadata Metadata
	Title    string
	Path     string
	Sections []Section
	Raw      string
}

// Section represents a markdown heading and its nested content.
type Section struct {
	Level    int
	Heading  string
	Paragraphs []string
	Tables   []Table
	Lists    []List
}

// Table represents a markdown pipe table.
type Table struct {
	Headers []string
	Rows    [][]string
}

// List represents an ordered or unordered markdown list block.
type List struct {
	Ordered bool
	Items   []string
}

// ParseMarkdown parses headings, tables, lists, and paragraphs from policy content.
func ParseMarkdown(content string) ([]Section, error) {
	body := stripFrontMatter(content)
	lines := strings.Split(body, "\n")

	var sections []Section
	var current *Section
	var paragraphLines []string
	var currentList *List

	flushParagraph := func() {
		if current == nil || len(paragraphLines) == 0 {
			paragraphLines = nil
			return
		}
		text := strings.TrimSpace(strings.Join(paragraphLines, "\n"))
		if text != "" {
			current.Paragraphs = append(current.Paragraphs, text)
		}
		paragraphLines = nil
	}

	flushList := func() {
		if current == nil || currentList == nil || len(currentList.Items) == 0 {
			currentList = nil
			return
		}
		current.Lists = append(current.Lists, *currentList)
		currentList = nil
	}

	ensureSection := func(level int, heading string) {
		flushParagraph()
		flushList()

		if heading == "" {
			if current == nil {
				current = &Section{Level: 1, Heading: "Document Body"}
				sections = append(sections, *current)
				current = &sections[len(sections)-1]
			}
			return
		}

		section := Section{
			Level:   level,
			Heading: heading,
		}
		sections = append(sections, section)
		current = &sections[len(sections)-1]
	}

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			flushParagraph()
			flushList()
			i++
			continue
		}

		if matches := headingPattern.FindStringSubmatch(trimmed); len(matches) == 3 {
			ensureSection(len(matches[1]), matches[2])
			i++
			continue
		}

		if isTableRow(trimmed) && i+1 < len(lines) && isTableSeparator(strings.TrimSpace(lines[i+1])) {
			flushParagraph()
			flushList()

			if current == nil {
				ensureSection(1, "Document Body")
			}

			table, consumed := parseTable(lines[i:])
			current.Tables = append(current.Tables, table)
			i += consumed
			continue
		}

		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushParagraph()

			if current == nil {
				ensureSection(1, "Document Body")
			}

			if currentList == nil || currentList.Ordered {
				flushList()
				currentList = &List{Ordered: false}
			}

			currentList.Items = append(currentList.Items, strings.TrimSpace(trimmed[2:]))
			i++
			continue
		}

		if matches := orderedListPattern.FindStringSubmatch(trimmed); len(matches) == 2 {
			flushParagraph()

			if current == nil {
				ensureSection(1, "Document Body")
			}

			if currentList == nil || !currentList.Ordered {
				flushList()
				currentList = &List{Ordered: true}
			}

			currentList.Items = append(currentList.Items, matches[1])
			i++
			continue
		}

		flushList()
		if current == nil {
			ensureSection(1, "Document Body")
		}
		paragraphLines = append(paragraphLines, trimmed)
		i++
	}

	flushParagraph()
	flushList()

	return sections, nil
}

func stripFrontMatter(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if i+1 < len(lines) {
				return strings.Join(lines[i+1:], "\n")
			}
			return ""
		}
	}
	return content
}

func isTableRow(line string) bool {
	return strings.Count(line, "|") >= 2
}

func isTableSeparator(line string) bool {
	if !strings.Contains(line, "|") {
		return false
	}
	cells := splitTableRow(line)
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		if cell == "" {
			continue
		}
		for _, r := range cell {
			if r != '-' && r != ':' {
				return false
			}
		}
	}
	return true
}

func parseTable(lines []string) (Table, int) {
	headerCells := splitTableRow(strings.TrimSpace(lines[0]))
	table := Table{Headers: headerCells}

	consumed := 2
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !isTableRow(line) {
			break
		}
		table.Rows = append(table.Rows, splitTableRow(line))
		consumed++
	}

	return table, consumed
}

func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "|")
	if line == "" {
		return nil
	}

	rawCells := strings.Split(line, "|")
	cells := make([]string, 0, len(rawCells))
	for _, cell := range rawCells {
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells
}
