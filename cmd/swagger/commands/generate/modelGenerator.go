package generate

import (
	"binh-swagger/cmd/swagger/commands/helpers"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

const (
	filePermOwnerReadWrite = 0o600
)

func Model(config spec.ModelSpec) error {
	var buf bytes.Buffer
	tmpl, err := loadModelTemplate()
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, config)
	if err != nil {
		return err
	}

	outputFile := helpers.GetAbsoluteSanitiseFilePath(config.OutputPath, config.OutputFile)
	shouldReturn, err := helpers.EnsureOutputDirectoryExists(config.OutputPath, os.Stdout, os.Stdin)
	if shouldReturn {
		return err
	}

	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}

func loadModelTemplate() (*template.Template, error) {
	currentDir, err := getCurrentFileDir()
	if err != nil {
		return nil, err
	}

	templatePath := helpers.GetAbsoluteSanitiseFilePath(filepath.Join(currentDir, "templates"), "model_template.tmpl")
	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	return template.New("model").Parse(string(tmpl))
}

// todo add tests
// todo move this to helpers package.
func getCurrentFileDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to create directory")
	}
	return filepath.Dir(file), nil
}

// todo use this when creating the output path and saving files
// func executableDir() (string, error) {
// 	exePath, err := os.Executable()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Dir(exePath), nil
// }
