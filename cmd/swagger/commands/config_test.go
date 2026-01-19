package commands

import (
	"bytes"
	"fmt"
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

func Test_ParseFile(t *testing.T) {
	filepath := fmt.Sprintf("%s/%s", validFilePath, validFilename)
	f, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("Unexpected error when opening test file: %v", err)
	}

	cfg, err := parseFile[Config](f)
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	if cfg.Name != "sample-test" {
		t.Errorf("expected Name 'sample-test', got %q", cfg.Name)
	}

	if cfg.Version != 1 {
		t.Errorf("expected Version 1, got %d", cfg.Version)
	}

	if !cfg.Enabled {
		t.Errorf("expected Enabled true")
	}

	if cfg.Metadata.Author != "test-user" {
		t.Errorf("expected author 'test-user', got %q", cfg.Metadata.Author)
	}

	if len(cfg.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(cfg.Items))
	}

	if cfg.Items[0].Value != "alpha" {
		t.Errorf("expected first item value 'alpha', got %q", cfg.Items[0].Value)
	}
}

func Test_ValidateConfigFlags(t *testing.T) {
	cmd := &ConfigCommand{API: true, CMD: true}
	err := validateConfigFlags(cmd)
	if err == nil {
		t.Fatalf("Expected error for both flags set, got nil")
	}

	cmd = &ConfigCommand{API: false, CMD: false}
	err = validateConfigFlags(cmd)
	if err == nil {
		t.Fatalf("Expected error for no flags set, got nil")
	}

	cmd = &ConfigCommand{API: true, CMD: false}
	err = validateConfigFlags(cmd)
	if err != nil {
		t.Fatalf("Did not expect error for valid flag set, got: %v", err)
	}

	cmd = &ConfigCommand{API: false, CMD: true}
	err = validateConfigFlags(cmd)
	if err != nil {
		t.Fatalf("Did not expect error for valid flag set, got: %v", err)
	}
}

func setupConfigTests(filename string, filePath string) (*ConfigCommand, *bytes.Buffer) {
	var buff bytes.Buffer
	baseCommand := SetupBaseCommand(&buff)
	args := ConfigArgs{Filepath: filePath, Filename: filename}
	cmd := &ConfigCommand{
		BaseCommand: *baseCommand,
		API:         true,
		Args:        args,
	}

	return cmd, &buff
}
