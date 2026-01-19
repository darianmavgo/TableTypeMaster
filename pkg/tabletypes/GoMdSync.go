package tabletypes

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// ParseReadme reads the README.md content and extracts the TableType rows.
func ParseReadme(content string) ([]TableType, error) {
	var rows []TableType
	lines := strings.Split(content, "\n")

	inTable := false
	headerSeen := false
	separatorSeen := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") {
			if !inTable {
				inTable = true
			}
			if !headerSeen {
				headerSeen = true
				continue
			}
			if !separatorSeen {
				// Check if it's a separator line (e.g. |---|---|)
				if strings.Contains(trimmed, "---") {
					separatorSeen = true
					continue
				}
			}

			// Parse row
			parts := strings.Split(trimmed, "|")
			// Expected format: | name | file_ext | ... | list_column_types |
			// Split results in empty string at start and end if line starts/ends with |
			if len(parts) < 9 { // 7 columns + 2 outer empties = 9
				continue
			}

			// Clean whitespace
			row := TableType{
				Name:            strings.TrimSpace(parts[1]),
				FileExt:         strings.TrimSpace(parts[2]),
				MimeType:        strings.TrimSpace(parts[3]),
				MagicNumber:     strings.TrimSpace(parts[4]),
				ListTables:      strings.TrimSpace(parts[5]),
				ListColumns:     strings.TrimSpace(parts[6]),
				ListColumnTypes: strings.TrimSpace(parts[7]),
			}
			rows = append(rows, row)
		} else if inTable {
			// End of table
			break
		}
	}
	return rows, nil
}

// GenerateReadmeTable generates the markdown table string from the data.
func GenerateReadmeTable(data []TableType) string {
	var sb strings.Builder
	sb.WriteString("| name | file_ext | mimetype | magic_number | list_tables | list_columns | list_column_types |\n")
	sb.WriteString("|---|---|---|---|---|---|---|\n")
	for _, row := range data {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
			row.Name, row.FileExt, row.MimeType, row.MagicNumber, row.ListTables, row.ListColumns, row.ListColumnTypes))
	}
	return sb.String()
}

// GenerateGoDataContent generates the content for data.go.
func GenerateGoDataContent(data []TableType) (string, error) {
	const tmpl = `package tabletypes

var AllTableTypes = []TableType{
{{- range . }}
	{Name: "{{.Name}}", FileExt: "{{.FileExt}}", MimeType: "{{.MimeType}}", MagicNumber: "{{.MagicNumber}}", ListTables: "{{.ListTables}}", ListColumns: "{{.ListColumns}}", ListColumnTypes: "{{.ListColumnTypes}}"},
{{- end }}
}
`
	t := template.Must(template.New("data").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Sync merges two slices of TableType.
// It returns a unified slice containing all unique entries from both inputs.
// If an entry exists in both, the one from 'primary' (list1) is preferred if they differ?
// The requirement is "catch the other table up". We'll assume strict union by Name.
// If Name matches, we assume they are the same or keep the one from list1.
func Sync(list1, list2 []TableType) []TableType {
	seen := make(map[string]TableType)
	var keys []string

	// Add list1
	for _, item := range list1 {
		if _, exists := seen[item.Name]; !exists {
			seen[item.Name] = item
			keys = append(keys, item.Name)
		}
	}

	// Add list2 (only if not exists)
	for _, item := range list2 {
		if _, exists := seen[item.Name]; !exists {
			seen[item.Name] = item
			keys = append(keys, item.Name)
		}
	}

	// Sort keys to ensure deterministic order (e.g. by Name)
	// Original order might be preserved better if we just iterated, but merging requires some order strategy.
	// Let's sort by definition order of standard types maybe? Or just alphabetical?
	// The README was not alphabetical. It was roughly usage based.
	// To preserve order: We can append new keys to the end.
	// But "larger table catches up" implies we might insert.
	// Let's just append new ones from list2 to the end of list1's order.

	// Re-construct the list
	// Actually, the 'seen' approach above does exactly that: keeps list1 order, then appends new list2 items.

	result := make([]TableType, 0, len(keys))
	for _, key := range keys {
		result = append(result, seen[key])
	}
	return result
}

// UpdateReadmeFile reads the file, replaces the table, and writes it back.
func UpdateReadmeFile(filepath string, data []TableType) error {
	contentBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	// Regex to find the table. Assumes table starts with header.
	// | name | file_ext | ...
	tableRegex := regexp.MustCompile(`(?m)^\| name \|.*$`)
	loc := tableRegex.FindStringIndex(content)
	if loc == nil {
		return fmt.Errorf("could not find table start in README.md")
	}
	startIdx := loc[0]

	// Find end of table (first line not starting with | after start)
	endIdx := len(content)
	lines := strings.Split(content[startIdx:], "\n")
	offset := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			// empty line might be end
			break // Usually markdown tables end with newline
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			break
		}
		offset += len(line) + 1 // +1 for newline
	}
	endIdx = startIdx + offset
    if endIdx > len(content) {
        endIdx = len(content)
    }

	newTable := GenerateReadmeTable(data)

	// Reassemble
	newContent := content[:startIdx] + newTable + content[endIdx:]

	return os.WriteFile(filepath, []byte(newContent), 0644)
}
