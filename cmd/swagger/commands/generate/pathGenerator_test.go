package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"binh-swagger/cmd/swagger/commands/internal/testhelpers"
	"os"
	"strings"
	"testing"
)

const (
	fileModeExecutable = 0o755
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
	resetProjectStructreForTests()
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

func Test_HandlerOperation_ShouldCreateHandlerFile(t *testing.T) {
	resetProjectStructreForTests()
	tempDir := t.TempDir()
	path := &spec.PathSpec{
		Name: "testPath",
		Get: &spec.Operation{
			Summary:     "Test GET operation",
			OperationID: "getTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
	}

	mockFileHelper := testhelpers.CreateMockFileHelper()
	if err := SetProjectStructure(tempDir); err != nil {
		t.Fatalf("Failed to set project structure: %v", err)
	}
	testProjectStructure, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get project structure: %v", err)
	}

	handlerDir := testProjectStructure["handlers"]
	if err = os.MkdirAll(handlerDir, fileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err = handleOperation("GET", path.Name, path.Get, mockFileHelper)
	if err != nil {
		t.Errorf("Expected no error for valid operation, got: %v", err)
	}

	expectedFilePath := mockFileHelper.GetAbsoluteSanitiseFilePath(
		tempDir+"/internal/handlers",
		path.Name+handlerFileSuffix,
	)

	_, err = os.Stat(expectedFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Expected handler file to be created at %s, but file does not exist", expectedFilePath)
			return
		}
		t.Errorf("Expected handler file to be created at %s, but got error: %v", expectedFilePath, err)
	}
}
