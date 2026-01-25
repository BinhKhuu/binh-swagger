package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"binh-swagger/cmd/swagger/commands/internal/testhelpers"
	"strings"
	"testing"
)

func Test_HnaldeOperation_ShouldThrowErrorWhenProjectNotSetup(t *testing.T) {
	path := &spec.PathSpec{
		Name: "testPath",
		Get:  &spec.Operation{},
	}
	mockFileHelper := testhelpers.CreateMockFileHelper()
	err := handleOperation("GET", path.Name, path.Get, mockFileHelper)
	if err == nil {
		t.Error("Expected error when project structure not initialized, got nil")
	}

	strings.Contains(err.Error(), "project structure not initialized")
}

func Test_HandleOperation_ShouldIgnoreOmittedOperations(t *testing.T) {
	tempDir := t.TempDir()
	path := &spec.PathSpec{
		Name: "testPath",
	}
	mockFileHelper := testhelpers.CreateMockFileHelper()
	err := SetProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to set project structure: %v", err)
	}
	err = handleOperation("GET", path.Name, path.Get, mockFileHelper)
	if err != nil {
		t.Errorf("Expected no error for omitted operation, got: %v", err)
	}
}
