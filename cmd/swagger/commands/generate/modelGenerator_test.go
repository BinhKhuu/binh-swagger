package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"binh-swagger/cmd/swagger/commands/internal/testhelpers"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func Test_GenerateModel(t *testing.T) {
	_ = InitGenerateTests(t)
	config := spec.ModelSpec{
		Name: "User",
		Fields: []spec.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}
	modelCmd, err := SpecToModelCommand(config)
	if err != nil {
		t.Fatalf("SpecToModelCommand failed: %v", err)
	}
	helperAdapter := testhelpers.CreateMockFileHelper()
	generateConfig := &testhelpers.MockGenerateConfig{
		FileHelperFunc: func() filehelper.FileHelper {
			return helperAdapter
		},
	}

	err = Model(modelCmd, generateConfig)
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
