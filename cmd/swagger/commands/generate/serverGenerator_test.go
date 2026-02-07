package generate

import (
	"os"
	"strings"
	"testing"
)

func Test_Server_Success(t *testing.T) {
	_ = InitGenerateTests(t)
	createMockTempFolders(t)
	cmd := &ServerCommand{
		PathNames: []string{"testPath1", "testPath2"},
	}
	cfg := createMockConfig()
	err := Server(cmd, cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	f, err := os.ReadFile(projectStructure["server"] + "/server.go")
	if err != nil {
		t.Fatalf("Expected server.go to be created, got error: %v", err)
	}

	contentStr := string(f)
	if !strings.Contains(contentStr, "package main") {
		t.Errorf("Expected server.go to contain 'package main', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "func main()") {
		t.Errorf("Expected server.go to contain 'func main()', got: %s", contentStr)
	}
}
