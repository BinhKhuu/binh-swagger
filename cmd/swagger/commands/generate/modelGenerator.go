package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	filePermOwnerReadWrite = 0o600
)

func Model(cmd *ModelCommand, generateConfig GenerateConfig) error {
	var buf bytes.Buffer
	fHelper := generateConfig.FileHelper()
	tmpl, err := loadModelTemplate(fHelper)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, cmd)
	if err != nil {
		return err
	}

	// todo think about coupling to GetProjectStructure here
	projectStructure, err := GetProjectStructure()
	if err != nil {
		return err
	}
	outputPath := projectStructure["models"]
	outFile := strings.ToLower(cmd.Name) + ".go"

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(outputPath, outFile)
	shouldReturn, err := fHelper.EnsureOutputDirectoryExists(outputPath, os.Stdout, os.Stdin)
	if shouldReturn {
		return err
	}

	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}

func loadModelTemplate(fileHefileHelper fileHelper.FileHelper) (*template.Template, error) {
	currentDir, err := getCurrentFileDir()
	if err != nil {
		return nil, err
	}

	templatePath := fileHefileHelper.GetAbsoluteSanitiseFilePath(filepath.Join(currentDir, "templates"), "model_template.tmpl")
	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	return template.New("model").Parse(string(tmpl))
}

// getCurrentFileDir returns the directory of the current file needs to be here to get templates based on relative path.
func getCurrentFileDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to create directory")
	}
	return filepath.Dir(file), nil
}
