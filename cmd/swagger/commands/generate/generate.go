package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
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
	ModelImportPath string
	Name            string
	Get             *spec.Operation
	Post            *spec.Operation
	Put             *spec.Operation
	Patch           *spec.Operation
	Delete          *spec.Operation
}

func SpecToPathCommand(cfg spec.PathSpec, pathName string) (*PathCommand, error) {
	paths, err := GetProjectStructure()
	if err != nil {
		return nil, err
	}
	return &PathCommand{
		ModelImportPath: projectImportPath + "/" + paths["models"],
		Name:            strings.Replace(pathName, "/", "", 1),
		Get:             cfg.Get,
		Post:            cfg.Post,
		Put:             cfg.Put,
		Patch:           cfg.Patch,
		Delete:          cfg.Delete,
	}, nil
}

func SpecToModelCommand(cfg spec.ModelSpec) (*ModelCommand, error) {
	return &ModelCommand{
		Name:   cfg.Name,
		Fields: cfg.Fields,
	}, nil
}
