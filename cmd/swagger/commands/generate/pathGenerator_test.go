package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"
	"unicode"
)

func createImportData() templateHelper.ImportsTemplateModel {
	return templateHelper.ImportsTemplateModel{
		ModelImportPath:   "github.com/example/project/models",
		HandlerImportPath: "github.com/example/project/handlers",
	}
}

func createMockOperationModel(cmd *PathCommand) []templateHelper.OperationModel {
	ops := createOperationModels(cmd)
	return ops
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

func createPathTestMocks() (*PathCommand, Config, templateHelper.ImportsTemplateModel) {
	return createMockPathCmd(), createMockConfig(), createImportData()
}

func createMockTempFolders(t *testing.T) {
	var err error
	// Create the Handlers and Routes temporary directories
	handlerDir := projectStructure["handlers"]
	if err = os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}
	routesDir := projectStructure["routes"]
	if err = os.MkdirAll(routesDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}
	serverDir := projectStructure["server"]
	if err = os.MkdirAll(serverDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}
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

func Test_createOperationModels_ShouldReturnCorrectOperations(t *testing.T) {
	pathCommand, _, _ := createPathTestMocks()

	ops := createOperationModels(pathCommand)
	if len(ops) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(ops))
	}
	for _, op := range ops {

		r := []rune(op.OperationID)[0]
		if !unicode.IsUpper(r) {
			t.Errorf("Expedted operation ID to be in pascal case but first letter was not uppercase: %s", op.OperationID)
		}
		if strings.ContainsAny(op.OperationID, " \t\n\r") {
			t.Errorf("Expected operation ID to be in pascal case but found whitespace: %s", op.OperationID)
		}
		if op.MethodType == "DELETE" {
			t.Error("Did not expect DELETE operation, but found in operations list")
		}
		if op.MethodType == "PATCH" {
			t.Error("Did not expect PATCH operation, but found in operations list")
		}
	}
}

func Test_Path_ShouldThrowErrorWhenProjectNotSetup(t *testing.T) {
	resetProjectStructreForTests()
	pathCommand, config, _ := createPathTestMocks()
	err := Path(pathCommand, config)
	if err != nil {
		if err.Error() != ErrProjectNotInitilized.Error() {
			t.Errorf("Expected error message to contain '%s', got: %v", ErrProjectNotInitilized, err)
		}
	} else {
		t.Error("Expected error when project structure is not set up, got nil")
	}
}

func Test_Path_ShouldGenerateHandlersAndRoutes(t *testing.T) {
	resetProjectStructreForTests()
	InitGenerateTests(t)
	createMockTempFolders(t)
	pathCommand, config, _ := createPathTestMocks()
	err := Path(pathCommand, config)
	if err != nil {
		t.Errorf("Expected no error when generating path, got: %v", err)
	}
}

func Test_CreateHandlers_ShouldCreateHandlerFile(t *testing.T) {
	InitGenerateTests(t)
	createMockTempFolders(t)
	path, _ := createMockSpecAndOperation()
	pathCommand, config, importData := createPathTestMocks()
	testProjectStructure, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get project structure: %v", err)
	}

	handlerDir := testProjectStructure["handlers"]

	err = createHandlers(pathCommand, config, importData, createMockOperationModel(pathCommand))
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
	createMockSpecAndOperation()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	if err := os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err := createHandlers(pathCommand, config, importData, createMockOperationModel(pathCommand))
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

	mockFileHelper := mockFileHelper.CreateMockFileHelper()
	testProjectStructure, _ := GetProjectStructure()
	handlerDir := testProjectStructure["handlers"]
	if err := os.MkdirAll(handlerDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}

	err := createHandlers(pathCommand, config, importData, createMockOperationModel(pathCommand))
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

	opts := []templateHelper.OperationModel{
		toOperationsModel(*path.Get, "GET"),
		toOperationsModel(*path.Post, "POST"),
		toOperationsModel(*path.Put, "PUT"),
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
	if route.OperationID != "getTestPath" {
		t.Errorf("Expected operation ID 'getTestPath', got: %s", route.OperationID)
	}
}

func Test_GetBaseTemplate_LoadsTemplateSuccessfully(t *testing.T) {
	testCases := []struct {
		templateName  string
		errorExpected bool
	}{
		{"route", false},
		{"handler", false},
		{"nonexistent", true},
	}

	for _, tc := range testCases {
		t.Run(tc.templateName, func(t *testing.T) {
			resetProjectStructreForTests()
			buf := bytes.Buffer{}
			InitGenerateTests(t)
			_, config, _ := createPathTestMocks()

			tmpl, err := templateHelper.GetBaseTemplate(&buf, config.FileHelper(), tc.templateName)

			if tc.errorExpected {
				if err == nil {
					t.Errorf("Expected error loading template '%s', got nil", tc.templateName)
				}
				if tmpl != nil {
					t.Error("Expected template to be nil on error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error loading template, got: %v", err)
				}
				if tmpl == nil {
					t.Error("Expected template to be non-nil")
				}
			}
		})
	}
}
