package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	validFilePath    = "testdata"
	validFilename    = "testconfig.yaml"
	invalidFilename  = "ugh.txt"
	invalidDirectory = "testdirectory"
)

func Test_GetSantisedFilePath(t *testing.T) {
	expectedAbs, _ := filepath.Abs(filepath.Join("..", validFilePath, validFilename))

	tests := []struct {
		filelocation string
		filename     string
		expected     string
	}{
		{filepath.Join("..", validFilePath), validFilename, expectedAbs},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.filelocation, tt.filename), func(t *testing.T) {
			result := GetSanitiseFilePath(tt.filelocation, tt.filename)
			if result != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, result)
			}
		})
	}
}

func Test_ValidateFileInfo(t *testing.T) {
	sanitisedFilePath, _ := filepath.Abs(filepath.Join("..", validFilePath, validFilename))
	fileInfo, err := os.Stat(sanitisedFilePath)
	if err != nil {
		t.Fatalf("Failed to stat file before running real tests: %v", err)
	}

	err = ValidateFileInfo(fileInfo)
	if err != nil {
		t.Errorf("Expected no error for valid file, but got: %v", err)
	}
}

func Test_DirectoryFilePath(t *testing.T) {
	sanitisedFilePath, _ := filepath.Abs(filepath.Join("..", validFilePath))
	fileInfo, err := os.Stat(sanitisedFilePath)
	if err != nil {
		t.Fatalf("Failed to stat file before running real tests: %v", err)
	}

	err = ValidateFileInfo(fileInfo)
	if err == nil {
		t.Fatalf("Expected Error but got nil")
	}

	if !strings.Contains(err.Error(), "directory") {
		t.Fatalf("expected directory error, got: %v", err)
	}
}
