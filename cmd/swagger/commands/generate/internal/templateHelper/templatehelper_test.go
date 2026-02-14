package templatehelper

import (
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"errors"
	"slices"
	"strings"
	"testing"
)

func Test_GetTemplatePath(t *testing.T) {
	tests := []struct {
		name         string
		templateKey  string
		expectedPath string
		expectError  bool
	}{
		{
			name:         "returns path for valid model template key",
			templateKey:  ModelTemplateKey,
			expectedPath: "model_template.tmpl",
			expectError:  false,
		},
		{
			name:         "returns path for valid handler template key",
			templateKey:  HandlerTemplateKey,
			expectedPath: "handler_template.tmpl",
			expectError:  false,
		},
		{
			name:        "returns error for non-existent template key",
			templateKey: "non_existent_template_key",
			expectError: true,
		},
		{
			name:        "returns error for empty template key",
			templateKey: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getTemplatePath(tt.templateKey)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error for template key %s, got nil", tt.templateKey)
				}
				if !errors.Is(err, ErrTemplateNotFound) {
					t.Fatalf("Expected ErrTemplateNotFound, got: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error for template key %s: %v", tt.templateKey, err)
				}
				if path != tt.expectedPath {
					t.Fatalf("Expected path %s, got %s", tt.expectedPath, path)
				}
			}
		})
	}
}

func Test_LoadModelTemplate(t *testing.T) {
	tests := []struct {
		name        string
		templateKey string
		expectError bool
	}{
		{
			name:        "loads model template successfully",
			templateKey: ModelTemplateKey,
			expectError: false,
		},
		{
			name:        "returns error for invalid template key",
			templateKey: "invalid_key",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileHelper := mockFileHelper.CreateMockFileHelper()
			tmpl, err := LoadModelTemplate(mockFileHelper, tt.templateKey)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error for template key %s, got nil", tt.templateKey)
				}
				if tmpl != nil {
					t.Fatalf("Expected nil template on error, got: %v", tmpl)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error loading template %s: %v", tt.templateKey, err)
				}
				if tmpl == nil {
					t.Fatalf("Expected valid template, got nil")
				}
			}
		})
	}
}

func Test_GetCurrentFileDir(t *testing.T) {
	dir, err := getCurrentFileDir()
	if err != nil {
		t.Fatalf("Unexpected error getting current file directory: %v", err)
	}
	if dir == "" {
		t.Fatalf("Expected non-empty directory path, got empty string")
	}
}

func Test_buildImportStringSlice_ReturnsImportPaths(t *testing.T) {
	inputPaths := map[string]string{
		"github.com/some/package": "github.com/some/package",
		"github.com/another/pkg":  "github.com/another/pkg",
	}
	expected := []string{"github.com/some/package", "github.com/another/pkg"}

	result := buildImportStringSlice(inputPaths)

	if len(result) != len(expected) {
		t.Fatalf("Expected %d import paths, got %d", len(expected), len(result))
	}

	for _, path := range expected {
		found := slices.Contains(result, path)
		if !found {
			t.Fatalf("Expected import path %s not found in result", path)
		}
	}
}

func Test_GetImportPath_ReturnsCorrectImportPath(t *testing.T) {
	tests := []struct {
		fieldType string
		expected  string
	}{
		{"time.Time", "time"},
		{"time.Duration", "time"},
		{"uuid.UUID", "github.com/google/uuid"},
		{"string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.fieldType, func(t *testing.T) {
			result := getImportPath(tt.fieldType)
			if result != tt.expected {
				t.Fatalf("Expected import path %s for field type %s, got %s", tt.expected, tt.fieldType, result)
			}
		})
	}
}

func Test_GetModelTemplateImportPaths_ReturnsPaths(t *testing.T) {
	tests := []struct {
		testName string
		expected []string
		tt       ModelTemplateModel
	}{
		{
			testName: "TestTimeGuidAndStringFields",
			expected: []string{
				"time",
				"github.com/google/uuid",
			},
			tt: ModelTemplateModel{
				Name: "TestModel",
				Fields: []FieldsModel{
					{Name: "CreatedAt", Type: "time.Time"},
					{Name: "ID", Type: "uuid.UUID"},
					{Name: "Name", Type: "string"},
				},
			},
		},
		{
			testName: "NoImportPathsForStringFields",
			expected: []string{},
			tt: ModelTemplateModel{
				Name: "TestModel",
				Fields: []FieldsModel{
					{Name: "ID", Type: "int"},
					{Name: "Name", Type: "string"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			SetModelTemplateImportPaths(&test.tt)

			expected := test.expected
			slices.Sort(expected)
			expectedStr := strings.Join(expected, ",")
			actual := test.tt.ImportPath
			slices.Sort(actual)

			actualStr := strings.Join(actual, ",")
			if expectedStr != actualStr {
				t.Fatalf("Expected import paths %s, got %s", expectedStr, actualStr)
			}
		})
	}
}
