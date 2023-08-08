package util

import (
	"os"
	"testing"
)

// TestExportCSV
func TestExportCSV(t *testing.T) {
	type Person struct {
		Name  string `csv:"Name"`
		Age   int    `csv:"Age"`
		Email string `csv:"Email"`
	}
	tbody := []Person{
		{Name: "Alice", Age: 25, Email: "alice@example.com"},
		{Name: "Bob", Age: 30, Email: "bob@example.com"},
	}
	buf, err := ExportCSV(Person{}, tbody, false)
	if err != nil {
		t.Errorf("export error: %s", err)
		return
	}
	os.WriteFile("test.csv", buf.Bytes(), 0644)
}
