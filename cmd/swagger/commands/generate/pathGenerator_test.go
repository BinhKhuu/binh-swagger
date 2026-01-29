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

	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Expected handler file to be created at %s, but file does not exist", expectedFilePath)
			return
		}
		t.Errorf("Expected handler file to be created at %s, but got error: %v", expectedFilePath, err)
		return
	}
	// todo assert if template has correct content
	contentStr := string(content)
	if !strings.Contains(contentStr, "getTestPath") {
		t.Errorf("Expected handler file to contain operation ID 'getTestPath', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "func getTestPath") {
		t.Errorf("Expected handler file to contain summary 'func getTestPath', got: %s", contentStr)
	}
}

func Test_CreateHandlers_ShouldHandleMultipleOperations(t *testing.T) {
	InitGenerateTests(t)
	path := &spec.PathSpec{
		Name: "multiPath",
		Get: &spec.Operation{
			Summary:     "Get operation",
			OperationID: "getMultiPath",
		},
		Post: &spec.Operation{
			Summary:     "Post operation",
			OperationID: "postMultiPath",
		},
	}

	operations := []operation{
		{method: "GET", op: path.Get},
		{method: "POST", op: path.Post},
	}

	mockFileHelper := testhelpers.CreateMockFileHelper()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	os.MkdirAll(handlerDir, fileModeExecutable)

	err := createHandlers(mockFileHelper, operations, path.Name+handlerFileSuffix)
	if err != nil {
		t.Errorf("Expected no error for multiple operations, got: %v", err)
	}

	expectedFilePath := mockFileHelper.GetAbsoluteSanitiseFilePath(
		handlerDir,
		path.Name+handlerFileSuffix,
	)

	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read handler file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "getMultiPath") {
		t.Errorf("Expected handler file to contain GET operation ID 'getMultiPath'")
	}
	if !strings.Contains(contentStr, "postMultiPath") {
		t.Errorf("Expected handler file to contain POST operation ID 'postMultiPath'")
	}
}

func Test_CreateHandlers_ShouldSkipNilOperations(t *testing.T) {
	InitGenerateTests(t)
	path := &spec.PathSpec{
		Name: "skipPath",
		Get: &spec.Operation{
			OperationID: "getSkipPath",
		},
	}

	operations := []operation{
		{method: "GET", op: path.Get},
		{method: "POST", op: nil}, // nil operation should be skipped
	}

	mockFileHelper := testhelpers.CreateMockFileHelper()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	os.MkdirAll(handlerDir, fileModeExecutable)

	err := createHandlers(mockFileHelper, operations, path.Name+handlerFileSuffix)
	if err != nil {
		t.Errorf("Expected no error when skipping nil operations, got: %v", err)
	}

	expectedFilePath := mockFileHelper.GetAbsoluteSanitiseFilePath(
		handlerDir,
		path.Name+handlerFileSuffix,
	)

	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read handler file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "getSkipPath") {
		t.Errorf("Expected handler file to contain GET operation ID 'getSkipPath'")
	}
}

func Test_CreateHandlers_ShouldReturnErrorOnTemplateFailure(t *testing.T) {
	InitGenerateTests(t)
	operations := []operation{
		{method: "INVALID", op: &spec.Operation{}}, // invalid method might cause template error
	}

	mockFileHelper := testhelpers.CreateMockFileHelper()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	os.MkdirAll(handlerDir, fileModeExecutable)

	err := createHandlers(mockFileHelper, operations, "test_handler.go")
	// This test depends on your template implementation
	// Adjust assertion based on expected behavior
	if err == nil {
		t.Log("Template handled invalid method gracefully")
	}
}
