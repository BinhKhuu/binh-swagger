package generate

import "testing"

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
