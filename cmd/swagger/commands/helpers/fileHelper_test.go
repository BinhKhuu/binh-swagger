package helpers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	validFilePath = "testdata"
	validFilename = "testconfig.yaml"
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
			result := GetAbsoluteSanitiseFilePath(tt.filelocation, tt.filename)
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

func Test_EnsureOutputDirectoryExists_CreateAndSucceeds(t *testing.T) {
	newDir := filepath.Join(t.TempDir(), "newdir")
	var outBuf bytes.Buffer
	input := strings.NewReader("y\n")

	shouldReturn, err := EnsureOutputDirectoryExists(newDir, &outBuf, input)
	if shouldReturn {
		t.Errorf("Expected shouldReturn to be false, but got true with error: %v", err)
	}
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if _, statErr := os.Stat(newDir); os.IsNotExist(statErr) {
		t.Error("Directory was not created after user said yes.")
	}
}

func Test_EnsureOutputDirectoryExists_CancelCreation(t *testing.T) {
	newDir := filepath.Join(t.TempDir(), "newdir")
	var outBuf bytes.Buffer
	input := strings.NewReader("n\n")

	shouldReturn, err := EnsureOutputDirectoryExists(newDir, &outBuf, input)
	if !shouldReturn {
		t.Errorf("Expected shouldReturn to be true, but got false with error: %v", err)
	}
	if err == nil {
		t.Error("Expected error when user cancels directory creation, but got nil")
	}
	if _, statErr := os.Stat(newDir); !os.IsNotExist(statErr) {
		t.Error("Directory was created despite user cancelling creation.")
	}
}

func Test_promptCreateDirectory_Yes(t *testing.T) {
	var outBuf bytes.Buffer
	input := strings.NewReader("y\n")

	result := promptCreateDirectory("somedir", &outBuf, input)
	if !result {
		t.Error("Expected promptCreateDirectory to return true for 'y' input, but got false")
	}
}

func Test_promptCreateDirectory_No(t *testing.T) {
	var outBuf bytes.Buffer
	input := strings.NewReader("n\n")

	result := promptCreateDirectory("somedir", &outBuf, input)
	if result {
		t.Error("Expected promptCreateDirectory to return false for 'n' input, but got true")
	}
}
