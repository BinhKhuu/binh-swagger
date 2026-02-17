package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"os"
	"testing"
)

func resetProjectStructreForTests() {
	projectStructure = nil
}

func InitGenerateTests(t *testing.T) string {
	tempDir := t.TempDir()
	resetProjectStructreForTests()
	err := SetProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to set project structure: %v", err)
	}
	return tempDir
}

func createMockData() (spec.Operation, map[string]*ModelCommand) {
	op := spec.Operation{
		Summary:     "Get user by ID",
		OperationID: "getUserByID",
		Produces:    []string{"application/json"},
		Responses: map[int]spec.ResponseSpec{
			200: {
				Description: "Successful response",
				Schema: &spec.SchemaSpec{
					Type: "object",
					Ref:  "#/definitions/User",
				},
			},
			400: {
				Description: "Bad request",
			},
		},
	}

	userModel := ModelCommand{
		Name: "User",
		Fields: []spec.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}

	mCommands := map[string]*ModelCommand{
		"User": &userModel,
	}

	return op, mCommands
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
	modelDir := projectStructure["models"]
	if err = os.MkdirAll(modelDir, pkg.FileModeExecutable); err != nil {
		t.Fatal(err)
	}
}

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
	models, _ := createMockModelCmd()
	modelsMap := make(map[string]*ModelCommand)
	modelsMap["TestModel"] = &models
	pathCmd, _ := SpecToPathCommand(*pathSpec, pathSpec.Name, modelsMap)
	return pathCmd
}

func createMockModelCmd() (ModelCommand, string) {
	model := ModelCommand{
		Name: "TestModel",
		Fields: []spec.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}
	key := "#/definitions/TestModel"
	return model, key
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

func createMockSpecAndOperation() (*spec.PathSpec, []operation) {
	_, key := createMockModelCmd()
	path := &spec.PathSpec{
		Name: "testPath",
		Get: &spec.Operation{
			Summary:     "Test GET operation",
			OperationID: "getTestPath",
			Produces:    []string{"application/json"},
			Responses: map[int]spec.ResponseSpec{
				200: {
					Description: "Successful response",
					Schema: &spec.SchemaSpec{
						Type: "object",
						Ref:  key,
					},
				},
			},
		},
		Post: &spec.Operation{
			Summary:     "Test POST operation",
			OperationID: "postTestPath",
			Produces:    []string{"application/json"},
			Responses: map[int]spec.ResponseSpec{
				201: {
					Description: "Created response",
					Schema:      &spec.SchemaSpec{},
				},
			},
		},
		Put: nil, // nil operation to test skipping
	}
	opts := []operation{
		{method: "GET", op: path.Get},
		{method: "POST", op: path.Post},
		{method: "PUT", op: path.Put}, // nil operation to test skipping
	}
	return path, opts
}
