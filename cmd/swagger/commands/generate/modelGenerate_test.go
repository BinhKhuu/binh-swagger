package generate

import (
	"binh-swagger/cmd/swagger/commands/helpers"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type MockFileHelper struct {
	ValidateFileInfoFn            func(os.FileInfo) error
	GetAbsoluteSanitiseFilePathFn func(string, string) string
	CheckSymlinksFn               func(string) error
	EnsureOutputDirectoryExistsFn func(string, io.Writer, io.Reader) (bool, error)
}

func (m MockFileHelper) ValidateFileInfo(fi os.FileInfo) error {
	return m.ValidateFileInfoFn(fi)
}

func (m MockFileHelper) GetAbsoluteSanitiseFilePath(loc, name string) string {
	return m.GetAbsoluteSanitiseFilePathFn(loc, name)
}

func (m MockFileHelper) CheckSymlinks(path string) error {
	return m.CheckSymlinksFn(path)
}

func (m MockFileHelper) EnsureOutputDirectoryExists(out string, w io.Writer, r io.Reader) (bool, error) {
	return m.EnsureOutputDirectoryExistsFn(out, w, r)
}

func mockFileHelper() MockFileHelper {
	mock := MockFileHelper{
		ValidateFileInfoFn: func(_ os.FileInfo) error {
			return nil
		},
		GetAbsoluteSanitiseFilePathFn: helpers.GetAbsoluteSanitiseFilePath,
		CheckSymlinksFn: func(_ string) error {
			return nil
		},
		EnsureOutputDirectoryExistsFn: func(_ string, _ io.Writer, _ io.Reader) (bool, error) {
			return true, nil
		},
	}

	return mock
}

func Test_GenerateModel(t *testing.T) {
	newDir := filepath.Join(t.TempDir(), "newdir")
	config := spec.ModelSpec{
		PackageName: "TestUser",
		Name:        "User",
		OutputPath:  newDir,
		OutputFile:  "models.go",
		Fields: []spec.FieldSpec{
			{Name: "ID", Type: "int", JSON: "id"},
			{Name: "Name", Type: "string", JSON: "name"},
		},
	}
	helperAdapter := mockFileHelper()
	// Call GenerateModel
	err := Model(config, helperAdapter)
	if err != nil {
		t.Fatalf("GenerateModel failed: %v", err)
	}

	// Verify output file exists
	outputDir := filepath.Join("..", "testdata/output")
	outputFile := filepath.Join(outputDir, "models.go")
	if _, err = os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file %s does not exist\n", outputFile)
	}

	// Optionally, read and verify contents of the output file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v\n", err)
	}

	expectedStrings := []string{
		"type User struct",
		"ID int",
		"Name string",
	}

	for _, str := range expectedStrings {
		if !bytes.Contains(content, []byte(str)) {
			t.Errorf("Expected string %q not found in output\n", str)
		}
	}
}
