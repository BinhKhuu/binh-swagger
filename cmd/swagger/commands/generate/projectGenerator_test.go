package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"path"
	"testing"
)

func Test_Project(t *testing.T) {
	tempDir := t.TempDir()
	rootDir := path.Join(tempDir, "myproject")

	fileHelper := &filehelper.DefaultFileHelper{} // testhelpers.CreateMockFileHelper()
	cfg := &spec.APIConfig{
		ProjectRoot: rootDir,
		Version:     "1.0.0",
		Models:      []spec.ModelSpec{},
		Routes:      []spec.RouteSpec{},
	}

	err := Project(cfg, fileHelper)
	if err != nil {
		t.Fatalf("Project generation failed: %v", err)
	}

	childDirs, err := fileHelper.ReadAllChildDirectoriesRecursive(rootDir)
	if err != nil {
		t.Fatalf("Failed to read child directories: %v", err)
	}

	expectedDirs := getProjectStructure(rootDir)

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
