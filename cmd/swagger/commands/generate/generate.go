package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"path/filepath"
	"strings"
)

type Config interface {
	FileHelper() filehelper.FileHelper
}

type ModelCommand struct {
	Name   string           `yaml:"name"`
	Fields []spec.FieldSpec `yaml:"fields"`
}

type PathCommand struct {
	ModelImportPath   string
	HandlerImportPath string
	Name              string
	Get               *spec.Operation
	Post              *spec.Operation
	Put               *spec.Operation
	Patch             *spec.Operation
	Delete            *spec.Operation
}

// todo refine the fields because this matches with serverTemplateModel
type ServerCommand struct {
	PathNames        []string
	RoutesImportPath string
}

// SpecToPathCommand todo test this.
func SpecToPathCommand(cfg spec.PathSpec, pathName string) (*PathCommand, error) {
	paths, err := GetProjectStructure()
	if err != nil {
		return nil, err
	}
	return &PathCommand{
		ModelImportPath:   filepath.Join(projectImportPath + "/" + paths["models"]),
		HandlerImportPath: filepath.Join(projectImportPath + "/" + paths["handlers"]),
		Name:              strings.Replace(pathName, "/", "", 1),
		Get:               cfg.Get,
		Post:              cfg.Post,
		Put:               cfg.Put,
		Patch:             cfg.Patch,
		Delete:            cfg.Delete,
	}, nil
}

// SpecToModelCommand todo test this.
func SpecToModelCommand(cfg spec.ModelSpec) (*ModelCommand, error) {
	return &ModelCommand{
		Name:   cfg.Name,
		Fields: cfg.Fields,
	}, nil
}
