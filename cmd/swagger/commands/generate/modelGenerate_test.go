package generate

import (
	"binh-swagger/cmd/swagger/commands"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func Test_GenerateModel(t *testing.T) {
	config := commands.ModelSpec{
		PackageName: "TestUser",
		Name:        "User",
		Fields: []commands.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}

	// Call GenerateModel
	err := GenerateModel(config)
	if err != nil {
		t.Fatalf("GenerateModel failed: %v", err)
	}

	// Verify output file exists
	outputDir := filepath.Join("..", "testdata/output")
	outputFile := filepath.Join(outputDir, "models.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s does not exist", outputFile)
	}

	// Optionally, read and verify contents of the output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedStrings := []string{
		"type User struct",
		"ID int",
		"Name string",
	}

	for _, str := range expectedStrings {
		if !bytes.Contains(content, []byte(str)) {
			t.Errorf("Expected string %q not found in output", str)
		}
	}
}
