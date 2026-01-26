package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"path"
	"testing"
)

func Test_Project(t *testing.T) {
	resetProjectStructreForTests()
	tempDir := t.TempDir()
	rootDir := path.Join(tempDir, "myproject")

	fileHelper := &filehelper.DefaultFileHelper{} // testhelpers.CreateMockFileHelper()
	cfg := &spec.APIConfig{
		ProjectRoot: rootDir,
		Version:     "1.0.0",
		Models:      map[string]spec.ModelSpec{},
		Paths:       map[string]spec.PathSpec{},
	}

	_, err := Project(cfg, fileHelper)
	if err != nil {
		t.Fatalf("Project generation failed: %v", err)
	}

	childDirs, err := fileHelper.ReadAllChildDirectoriesRecursive(rootDir)
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
