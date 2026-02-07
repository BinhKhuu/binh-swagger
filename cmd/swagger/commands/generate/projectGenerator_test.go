package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	mockFileHelper "binh-swagger/cmd/swagger/commands/helpers/mocks"
	"testing"
)

// Test_Project test relies on fileHelper real implmementation to create directories and read them back which is not good consider mocking and moving logic to filehelper.
func Test_Project(t *testing.T) {
	resetProjectStructreForTests()

	// we should mock this instead of creating the actual file consider refactoring the tests.
	tempDir := mockFileHelper.ChangeCWDAndCreateGoModFile(t)
	fileHelper := &filehelper.DefaultFileHelper{}

	_, err := Project(fileHelper, tempDir)
	if err != nil {
		t.Fatalf("Project generation failed: %v", err)
	}

	childDirs, err := fileHelper.ReadAllChildDirectoriesRecursive(tempDir)
	if err != nil {
		t.Fatalf("Failed to read child directories: %v", err)
	}

	expectedDirs, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get expected project structure: %v", err)
	}

	dirSet := make(map[string]bool)
	for _, dir := range childDirs {
		dirSet[dir] = true
	}

	// Verify all expected directories are created
	// Because of nested directories, we only check for expected directories returned by getProjectStructure
	for _, dir := range expectedDirs {
		if !dirSet[dir] {
			t.Errorf("Expected directory %s not found", dir)
		}
	}

	if projectImportPath != "testmodule" {
		t.Errorf("Expected project import path 'testmodule', got '%s'", projectImportPath)
	}
}

func Test_GetProjectStructure_InitReturnsDiretoryStructure(t *testing.T) {
	InitGenerateTests(t)
	structure, err := GetProjectStructure()
	if err != nil {
		t.Fatalf("Failed to get project structure: %v", err)
	}
	if structure == nil {
		t.Fatalf("Expected non-nil project structure, got nil")
	}

	expectedDirs := []string{"server", "handlers", "routes", "models"}
	for _, dirKey := range expectedDirs {
		if _, exists := structure[dirKey]; !exists {
			t.Errorf("Expected directory key %s not found in project structure", dirKey)
		}
	}
}

func Test_GetProjectStructure_NoInitThrowsError(t *testing.T) {
	resetProjectStructreForTests()
	_, err := GetProjectStructure()
	if err == nil {
		t.Fatalf("Expected error when project structure not initialized, got nil")
	}
}
