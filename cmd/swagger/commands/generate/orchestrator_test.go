package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	fileHelperMock "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"testing"
)

type mockConfig struct {
	fileHelper filehelper.FileHelper
}

func (m *mockConfig) FileHelper() filehelper.FileHelper {
	return m.fileHelper
}

func Test_FromAPIConfig_Success(t *testing.T) {
	resetProjectStructreForTests()
	_ = fileHelperMock.ChangeCWDAndCreateGoModFile(t)
	config := &mockConfig{
		fileHelper: &filehelper.DefaultFileHelper{},
	}
	apiConfig := createMockAPIConfig()
	if err := fromAPIConfig(apiConfig, config); err != nil {
		t.Fatalf("Error generating from API config: %v", err)
	}
}

func createMockAPIConfig() *spec.APIConfig {
	return &spec.APIConfig{
		Version: "1.0.0",
		Models: map[string]spec.ModelSpec{
			"User": {
				Name: "User",
				Fields: []spec.FieldSpec{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Name", Type: "string", JSON: "name"},
				},
			},
			"Product": {
				Name: "Product",
				Fields: []spec.FieldSpec{
					{Name: "ID", Type: "int", JSON: "id"},
					{Name: "Title", Type: "string", JSON: "title"},
					{Name: "Price", Type: "float64", JSON: "price"},
				},
			},
		},
		Paths: map[string]spec.PathSpec{
			"/users": {
				Get: &spec.Operation{
					Summary:     "Get all users",
					OperationID: "getUsers",
					Produces:    []string{"application/json"},
				},
			},
		},
	}
}
