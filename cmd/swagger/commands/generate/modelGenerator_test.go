package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func Test_GenerateModel(t *testing.T) {
	_ = InitGenerateTests(t)
	createMockTempFolders(t)
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
	helperAdapter := mockFileHelper.CreateMockFileHelper()
	generateConfig := &mockFileHelper.MockGenerateConfig{
		FileHelperFunc: func() filehelper.FileHelper {
			return helperAdapter
		},
	}

	err = Model(modelCmd, generateConfig)
	if err != nil {
		t.Fatalf("GenerateModel failed: %v", err)
	}

	outputFile := filepath.Join(projectStructure["models"], config.Name+".go")
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

func Test_ToModelTemplate_(t *testing.T) {
	tests := []struct {
		name     string
		cmd      ModelCommand
		expected templateHelper.ModelTemplateModel
	}{
		{
			name: "FieldsWithTimeType",
			cmd: ModelCommand{
				Name: "User",
				Fields: []spec.FieldSpec{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Name", Type: "string", JSON: "name"},
					{Name: "Timestamp", Type: "time.Time", JSON: "timestamp"},
					{Name: "", Type: "string", JSON: "invalid"},      // Invalid field missing Name
					{Name: "InvalidType", Type: "", JSON: "invalid"}, // Invalid field missing Type
				},
			},
			expected: templateHelper.ModelTemplateModel{
				Name: "User",
				Fields: []templateHelper.FieldsModel{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Name", Type: "string", JSON: "name"},
					{Name: "Timestamp", Type: "time.Time", JSON: "timestamp"},
				},
				ImportPath: []string{"time"},
			},
		},
		{
			name: "FieldsWithNoImportTypes",
			cmd: ModelCommand{
				Name: "User",
				Fields: []spec.FieldSpec{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Name", Type: "string", JSON: "name"},
					{Name: "", Type: "string", JSON: "invalid"},      // Invalid field missing Name
					{Name: "InvalidType", Type: "", JSON: "invalid"}, // Invalid field missing Type
				},
			},
			expected: templateHelper.ModelTemplateModel{
				Name: "User",
				Fields: []templateHelper.FieldsModel{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Name", Type: "string", JSON: "name"},
				},
				ImportPath: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modelTemplate := toModelTemplate(&test.cmd)

			if modelTemplate.Name != test.cmd.Name {
				t.Errorf("Expected model name %q, got %q\n", test.cmd.Name, modelTemplate.Name)
			}

			if len(modelTemplate.Fields) != len(test.expected.Fields) {
				t.Errorf("Expected %d valid fields, got %d\n", len(test.expected.Fields), len(modelTemplate.Fields))
			}

			if len(modelTemplate.ImportPath) != len(test.expected.ImportPath) {
				t.Errorf("Expected import paths %v, got %v\n", test.expected.ImportPath, modelTemplate.ImportPath)
			}
		})
	}
}
