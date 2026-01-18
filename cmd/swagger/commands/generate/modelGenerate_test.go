package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// todo use temp dir for output it will be cleaned up after test.
func Test_GenerateModel(t *testing.T) {
	config := spec.ModelSpec{
		PackageName: "TestUser",
		Name:        "User",
		OutputPath:  filepath.Join("..", "testdata", "output"),
		OutputFile:  "models.go",
		Fields: []spec.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}

	// Call GenerateModel
	err := Model(config)
	if err != nil {
		t.Fatalf("GenerateModel failed: %v", err)
	}

	// Verify output file exists
	outputDir := filepath.Join("..", "testdata/output")
	outputFile := filepath.Join(outputDir, "models.go")
	if _, err = os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s does not exist\n", outputFile)
	}

	// Optionally, read and verify contents of the output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v\n", err)
	}

	expectedStrings := []string{
		"type User struct",
		"ID int",
		"Name string",
	}

	for _, str := range expectedStrings {
		if !bytes.Contains(content, []byte(str)) {
			t.Errorf("Expected string %q not found in output\n", str)
		}
	}
}
