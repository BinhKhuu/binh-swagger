package generate

import (
	"binh-swagger/cmd/swagger/commands/helpers"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// filePermOwnerReadWrite is the file permission for owner read/write only (rw-------).
	filePermOwnerReadWrite = 0o600
	fileModeExecutable     = 0o755
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

	// Check if directory exists
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		if !promptCreateDirectory(config.OutputPath) {
			return errors.New("directory creation cancelled by user")
		}
		if err := os.MkdirAll(config.OutputPath, fileModeExecutable); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}

func promptCreateDirectory(path string) bool {
	reader := bufio.NewReader(os.Stdin)
	// todo somehow pass the output from the command to here instead of os.Stdout
	fmt.Fprintf(os.Stdout, "Directory %s does not exist. Create it? it will be made in the directory relative to where you ran the command [y/N]: ", path)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
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

// todo move this to helpers package
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
