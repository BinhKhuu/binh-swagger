package generate

import (
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"net/http"
	"os"
	"strings"
	"testing"
)

func createImportData() importsTemplateModel {
	return importsTemplateModel{
		ModelImportPath:   "github.com/example/project/models",
		HandlerImportPath: "github.com/example/project/handlers",
	}
}

func createMockPathCmd() *PathCommand {
	pathSpec, _ := createMockSpecAndOperation()
	return &PathCommand{
		Name: "testPath",
		Get:  pathSpec.Get,
		Post: pathSpec.Post,
	}
}

func createMockConfig() Config {
	mockFileHelper := mockFileHelper.CreateMockFileHelper()
	return &mockConfig{
		fileHelper: mockFileHelper,
	}
}

func createPathTestMocks() (*PathCommand, Config, importsTemplateModel) {
	return createMockPathCmd(), createMockConfig(), createImportData()
}

func createMockSpecAndOperation() (*spec.PathSpec, []operation) {
	path := &spec.PathSpec{
		Name: "testPath",
		Get: &spec.Operation{
			Summary:     "Test GET operation",
			OperationID: "getTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
		Post: &spec.Operation{
			Summary:     "Test POST operation",
			OperationID: "postTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
		Put: nil,
	}
	opts := []operation{
		{method: "GET", op: path.Get},
		{method: "POST", op: path.Post},
		{method: "PUT", op: path.Put}, // nil operation to test skipping
	}
	return path, opts
}

func Test_CreateHandlers_ShouldThrowErrorWhenProjectNotSetup(t *testing.T) {
	resetProjectStructreForTests()
	pathCommand, config, importData := createPathTestMocks()
	_, operations := createMockSpecAndOperation()

	err := createHandlers(pathCommand, config, importData, operations)
	if err == nil {
		t.Error("Expected error when project structure not initialized, got nil")
	}

	// if err is nil and we use err.Error() it will panic
	if err == nil || !strings.Contains(err.Error(), "project structure not initialized") {
		t.Errorf("Expected project structure not initialized error, got: %v", err)
	}
}

func Test_CreateHandlers_ShouldCreateHandlerFile(t *testing.T) {
	InitGenerateTests(t)
	path, operations := createMockSpecAndOperation()
	pathCommand, config, importData := createPathTestMocks()
	testProjectStructure, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get project structure: %v", err)
	}

	handlerDir := testProjectStructure["handlers"]
	if err = os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err = createHandlers(pathCommand, config, importData, operations)
	if err != nil {
		t.Errorf("Expected no error for valid operation, got: %v", err)
	}

	expectedFilePath := config.FileHelper().GetAbsoluteSanitiseFilePath(
		handlerDir,
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
	pathCommand, config, importData := createPathTestMocks()
	_, operations := createMockSpecAndOperation()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	if err := os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err := createHandlers(pathCommand, config, importData, operations)
	if err != nil {
		t.Errorf("Expected no error for multiple operations, got: %v", err)
	}

	expectedFilePath := config.FileHelper().GetAbsoluteSanitiseFilePath(
		handlerDir,
		pathCommand.Name+handlerFileSuffix,
	)

	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read handler file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "getTestPath") {
		t.Errorf("Expected handler file to contain GET operation ID 'getMultiPath'")
	}
	if !strings.Contains(contentStr, "postTestPath") {
		t.Errorf("Expected handler file to contain POST operation ID 'postMultiPath'")
	}
}

func Test_CreateHandlers_ShouldSkipNilOperations(t *testing.T) {
	InitGenerateTests(t)
	pathCommand, config, importData := createPathTestMocks()
	_, operations := createMockSpecAndOperation()

	mockFileHelper := mockFileHelper.CreateMockFileHelper()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	if err := os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err := createHandlers(pathCommand, config, importData, operations)
	if err != nil {
		t.Errorf("Expected no error when skipping nil operations, got: %v", err)
	}

	expectedFilePath := mockFileHelper.GetAbsoluteSanitiseFilePath(
		handlerDir,
		pathCommand.Name+handlerFileSuffix,
	)

	content, err := os.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read handler file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "PUT") {
		t.Errorf("Expected to skip NIL PUT operation, but found in handler file")
	}
}

func Test_CreateRoutesData_ReturnsRouteData(t *testing.T) {
	path := &spec.PathSpec{
		Name: "/testPath",
		Get: &spec.Operation{
			Summary:     "Test GET operation",
			OperationID: "getTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
		Post: &spec.Operation{
			Summary:     "Test POST operation",
			OperationID: "postTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
		Put: &spec.Operation{
			Summary:     "Test PUT operation",
			OperationID: "putTestPath",
			Produces:    []string{"application/json"},
			Responses:   map[int]spec.ResponseSpec{},
		},
	}

	opts := []operation{
		{method: "GET", op: path.Get},
		{method: "POST", op: path.Post},
		{method: "PUT", op: path.Put},
	}

	routes := createRoutesData(path.Name, opts)

	if len(routes) < 3 || len(routes) > 3 {
		t.Errorf("Expected 3 routes but go %d", len(routes))
	}

	// just asserting the first element
	route := routes[0]
	if route.Method != http.MethodGet {
		t.Errorf("Expected method 'GET', got: %s", route.Method)
	}
	if route.PathName != path.Name {
		t.Errorf("Expected path name '%s', got: %s", path.Name, route.PathName)
	}
	if route.OperationId != "getTestPath" {
		t.Errorf("Expected operation ID 'getTestPath', got: %s", route.OperationId)
	}
}

// func Test_GenerateRoutes(t *testing.T) {
// }

// // todo test imports routes
