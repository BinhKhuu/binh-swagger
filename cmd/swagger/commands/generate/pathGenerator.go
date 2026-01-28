package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"os"
)

const (
	handlerFileSuffix = "_handler.go"
)

func handleOperation(filename string, cfg *spec.Operation, fHelper fileHelper.FileHelper) error {
	if cfg == nil {
		return nil
	}
	_, err := GetProjectStructure()
	if err != nil {
		return err
	}

	// generate handler
	err = createHandlerFile(projectStructure["handlers"], filename, cfg, fHelper)
	if err != nil {
		return err
	}
	// generate route

	return nil
}

func createHandlerFile(handlersDir string, filename string, cfg *spec.Operation, fHelper fileHelper.FileHelper) error {
	var buf bytes.Buffer
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(handlersDir, filename+handlerFileSuffix)
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateHelper.HandlerTemplateKey)
	if err != nil {
		return err
	}

	if err = tmpl.Execute(&buf, cfg); err != nil {
		return err
	}

	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}

func Path(cmd *PathCommand, config Config) error {
	fHelper := config.FileHelper()
	ops := []struct {
		method string
		op     *spec.Operation
	}{
		{"GET", cmd.Get},
		{"POST", cmd.Post},
		{"PUT", cmd.Put},
		{"DELETE", cmd.Delete},
		{"PATCH", cmd.Patch},
	}

	for _, o := range ops {
		err := handleOperation(cmd.Name, o.op, fHelper)
		if err != nil {
			return err
		}
	}

	return nil
}
