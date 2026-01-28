package templateHelper

import (
	"binh-swagger/cmd/swagger/commands/internal/testhelpers"
	"errors"
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
			mockFileHelper := testhelpers.CreateMockFileHelper()
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
