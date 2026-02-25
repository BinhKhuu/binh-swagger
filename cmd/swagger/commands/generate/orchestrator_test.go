package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	fileHelperMock "binh-swagger/cmd/swagger/commands/helpers/mocks"
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
