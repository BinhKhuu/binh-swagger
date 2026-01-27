package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

// todo test this file.
var templatePaths = map[string]string{
	"model":  "model_template.tmpl",
	"hander": "handler_template.tmpl",
}

func LoadModelTemplate(fileHefileHelper fileHelper.FileHelper, templateName string) (*template.Template, error) {
	templateFilename, err := getTemplatePath(templateName)
	if err != nil {
		return nil, err
	}

	currentDir, err := getCurrentFileDir()
	if err != nil {
		return nil, err
	}

	templatePath := fileHefileHelper.GetAbsoluteSanitiseFilePath(filepath.Join(currentDir, "templates"), templateFilename)
	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	return template.New("model").Parse(string(tmpl))
}

func getTemplatePath(templateName string) (string, error) {
	path := templatePaths[templateName]
	if path == "" {
		return "", errors.New(templateName + " template not found")
	}
	return path, nil
}

// getCurrentFileDir returns the directory of the current file needs to be here to get templates based on relative path.
func getCurrentFileDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to create directory")
	}
	return filepath.Dir(file), nil
}
