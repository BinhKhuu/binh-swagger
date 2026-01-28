package templateHelper

import (
	"binh-swagger/cmd/swagger/commands/internal/testhelpers"
	"errors"
	"testing"
)

func Test_GetTemplatePath_ReturnsNotFound(t *testing.T) {
	badKey := "non_existent_template_key"
	_, err := getTemplatePath(badKey)
	if err == nil {
		t.Fatalf("Expected error for non-existent template key, got nil")
	}

	if err != ErrTemplateNotFound && !errors.Is(err, ErrTemplateNotFound) {
		t.Fatalf("Expected ErrTemplateNotFound, got: %v", err)
	}
}

func Test_GetTemplatePath_ReturnsPath(t *testing.T) {
	path, err := getTemplatePath(ModelTemplateKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedPath := "model_template.tmpl"
	if path != expectedPath {
		t.Fatalf("Expected path %s, got %s", expectedPath, path)
	}
}

func Test_LoadModelTemp(t *testing.T) {
	mockFileHelper := testhelpers.CreateMockFileHelper()
	_, err := LoadModelTemplate(mockFileHelper, ModelTemplateKey)
	if err != nil {
		t.Fatalf("Failed to load model template: %v", err)
	}
}
