package tabletypes

import (
	"os"
	"testing"
)

func TestExportToHTML(t *testing.T) {
	filename := "test_output.html"
	defer os.Remove(filename)

	err := ExportToHTML(filename)
	if err != nil {
		t.Fatalf("ExportToHTML failed: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("File %s was not created", filename)
	}
}

func TestExportToSQLite(t *testing.T) {
	filename := "test_output.sqlite"
	defer os.Remove(filename)

	err := ExportToSQLite(filename)
	if err != nil {
		t.Fatalf("ExportToSQLite failed: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("File %s was not created", filename)
	}
}
