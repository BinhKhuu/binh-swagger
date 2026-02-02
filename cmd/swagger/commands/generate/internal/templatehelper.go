package templatehelper

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
)

var ErrTemplateNotFound = errors.New("template not found")

var templatePaths = map[string]string{
	ModelTemplateKey:   "model_template.tmpl",
	HandlerTemplateKey: "handler_template.tmpl",
	RouteTemplateKey:   "route_template.tmpl",
}

const (
	HandlerTemplateKey = "handler"
	ModelTemplateKey   = "model"
	RouteTemplateKey   = "route"
)

// todo turn this into a struct or dict so you fetch like templatehelper.Templates.Routes.Body
const (
	RouteDefine   = "registerRoutes"
	ImportsDefine = "imports"
)

func LoadModelTemplate(fileHefileHelper fileHelper.FileHelper, templateName string) (*template.Template, error) {
	templateFilename, err := getTemplatePath(templateName)
	if err != nil {
		return nil, err
	}

	currentDir, err := getCurrentFileDir()
	if err != nil {
		return nil, err
	}

	// path of template is coupled to the location of this file.
	// also not properly tested unit test needs to reflect where this path is
	templatePath := fileHefileHelper.GetAbsoluteSanitiseFilePath(filepath.Join(currentDir, "..", "templates"), templateFilename)

	template := template.Must(template.ParseFiles(templatePath))
	return template, nil

	// old code here as an example of reading file content if needed
	// tmpl, err := os.ReadFile(templatePath)
	// if err != nil {
	// 	return nil, err
	// }
	// return template.New(templateName).Parse(string(tmpl))
}

func getTemplatePath(templateName string) (string, error) {
	path := templatePaths[templateName]
	if path == "" {
		return "", fmt.Errorf("%w: %s", ErrTemplateNotFound, templateName)
	}
	return path, nil
}

// getCurrentFileDir returns the directory of the current file needs to be here to get templates based on relative path.
func getCurrentFileDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get directory")
	}
	return filepath.Dir(file), nil
}
