package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"strings"
)

type GenerateConfig interface {
	FileHelper() filehelper.FileHelper
}

type GenerateCommand struct {
	Name   string
	Get    *spec.Operation
	Post   *spec.Operation
	Put    *spec.Operation
	Patch  *spec.Operation
	Delete *spec.Operation
}

func SpecToCommand(cfg spec.PathSpec, pathName string) (*GenerateCommand, error) {
	return &GenerateCommand{
		Name:   strings.Replace(pathName, "/", "", 1),
		Get:    cfg.Get,
		Post:   cfg.Post,
		Put:    cfg.Put,
		Patch:  cfg.Patch,
		Delete: cfg.Delete,
	}, nil
}
