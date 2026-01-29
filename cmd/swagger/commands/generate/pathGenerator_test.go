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

func Test_CreateHandlers_ShouldThrowErrorWhenProjectNotSetup(t *testing.T) {
	resetProjectStructreForTests()
	path, opt := createMockOperation()
	mockFileHelper := testhelpers.CreateMockFileHelper()
	operations := []operation{opt}
	err := createHandlers(mockFileHelper, operations, path.Name+handlerFileSuffix)
	if err == nil {
		t.Error("Expected error when project structure not initialized, got nil")
	}

	// if err is nil and we use err.Error() it will panic
	if err == nil || !strings.Contains(err.Error(), "project structure not initialized") {
		t.Errorf("Expected project structure not initialized error, got: %v", err)
	}
}

func createMockOperation() (*spec.PathSpec, operation) {
	path := &spec.PathSpec{
		Name: "testPath",
		Get:  &spec.Operation{},
	}
	operation := operation{
		method: "GET",
		op:     path.Get,
	}
	return path, operation
}

// no files written to need to assert of the buffer
func Test_CreateHandlers_ShouldCreateHandlerFile(t *testing.T) {
	tempDir := InitGenerateTests(t)
	path := &spec.PathSpec{
		Name: "testPath",
		Get: &spec.Operation{
			Summary:     "Test GET operation",
			OperationID: "getTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
	}
	opt := operation{
		method: "GET",
		op:     path.Get,
	}
	operations := []operation{opt}
	mockFileHelper := testhelpers.CreateMockFileHelper()
	testProjectStructure, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get project structure: %v", err)
	}

	handlerDir := testProjectStructure["handlers"]
	if err = os.MkdirAll(handlerDir, fileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err = createHandlers(mockFileHelper, operations, path.Name+handlerFileSuffix)
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

	// todo assert if template has correct content
}
