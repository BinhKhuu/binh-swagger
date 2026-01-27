package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"os"
)

const (
	handlerFileSuffix = "_handler.go"
)

func handleOperation(method string, filename string, cfg *spec.Operation, fHelper fileHelper.FileHelper) error {
	if cfg == nil {
		return nil
	}
	projectStructure, err := GetProjectStructure()
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
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(handlersDir, filename+handlerFileSuffix)
	// todo make tmpl file
	// parse spec into tmp file

	// save
	return os.WriteFile(outputFile, []byte("// Handler code here"), filePermOwnerReadWrite)
}

func Path(cmd *PathCommand, config GenerateConfig) error {
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
		err := handleOperation(o.method, cmd.Name, o.op, fHelper)
		if err != nil {
			return err
		}
	}

	return nil
}
