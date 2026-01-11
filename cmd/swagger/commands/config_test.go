package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

var (
	validFilePath    = "./testdata"
	validFilename    = "testconfig.yaml"
	invalidFilename  = "ugh.txt"
	invalidDirectory = "testdirectory"
)

func Test_ValidFilePath(t *testing.T) {
	cmd, _ := setupConfigTests(validFilename, validFilePath)

	err := cmd.Execute(nil)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
}

func Test_InvalidFilePath(t *testing.T) {
	cmd, _ := setupConfigTests(invalidFilename, validFilePath)
	err := cmd.Execute(nil)
	if err == nil {
		t.Fatalf("Expected Error but got nil")
	}

	if !os.IsNotExist(err) {
		t.Fatalf("Expected Is Not Exist error but got: %v", err)
	}
}

func setupConfigTests(filename string, filePath string) (*ConfigCommand, *bytes.Buffer) {
	var buff bytes.Buffer
	args := ConfigArgs{Filepath: filePath, Filename: filename}
	cmd := &ConfigCommand{
		Out:  &buff,
		Args: args,
	}

	return cmd, &buff
}

func Test_DirectoryFilePath(t *testing.T) {
	cmd, _ := setupConfigTests(invalidFilename, invalidDirectory)
	err := cmd.Execute(nil)
	if err == nil {
		t.Fatalf("Expected Error but got nil")
	}

	if !strings.Contains(err.Error(), "directory") {
		t.Fatalf("expected directory error, got: %v", err)
	}
}
